package module

import (
	"github.com/wudiliujie/common/chanrpc"
	"github.com/wudiliujie/common/console"
	"github.com/wudiliujie/common/go"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/timer"
	"time"
)

type IEngine interface {
	Update(diff int64)
}

type Skeleton struct {
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	ChanRPCServer      *chanrpc.Server
	g                  *g.Go
	dispatcher         *timer.Dispatcher
	client             *chanrpc.Client
	server             *chanrpc.Server
	commandServer      *chanrpc.Server
	strartTime         int64 //启动时间
	runTime            int64 //运行时间
	Name               string
	Engine             IEngine
	debugFunc          string
	debugData          interface{} //调试数据
	ticker             *time.Ticker
}

func (s *Skeleton) Init() {
	if s.GoLen <= 0 {
		s.GoLen = 0
	}
	if s.TimerDispatcherLen <= 0 {
		s.TimerDispatcherLen = 0
	}
	if s.AsynCallLen <= 0 {
		s.AsynCallLen = 0
	}

	s.g = g.New(s.GoLen)
	s.dispatcher = timer.NewDispatcher(s.TimerDispatcherLen)
	s.client = chanrpc.NewClient(s.AsynCallLen)
	s.server = s.ChanRPCServer

	if s.server == nil {
		s.server = chanrpc.NewServer(0)
	}
	s.commandServer = chanrpc.NewServer(0)

	s.strartTime = time.Now().UnixNano() / 1e6

	//s.AfterFunc(time.Microsecond*30, s.ThreadRun)
	s.ticker = time.NewTicker(time.Microsecond * 30)
}

func (s *Skeleton) Run(closeSig chan bool) {
	for {
		select {
		case <-closeSig:
			s.ticker.Stop()
			s.commandServer.Close()
			s.server.Close()
			for !s.g.Idle() || !s.client.Idle() {
				s.g.Close()
				s.client.Close()
			}
			return
		case <-s.ticker.C:
			s.ThreadRun()
		case ri := <-s.client.ChanAsynRet:
			s.debugData = ri
			s.debugFunc = "s.client.ChanAsynRet"
			s.client.Cb(ri)
		case ci := <-s.server.ChanCall:
			s.debugFunc = "s.server.ChanCall"
			s.debugData = s.server.Exec(ci).GetDebug()
		case ci := <-s.commandServer.ChanCall:
			s.debugData = ci
			s.debugFunc = "s.commandServer.ChanCall"
			s.commandServer.Exec(ci)
		case cb := <-s.g.ChanCb:
			s.debugFunc = "s.g.ChanCb"
			s.debugData = cb
			s.g.Cb(cb)
		case t := <-s.dispatcher.ChanTimer:
			s.debugData = "chantimer"
			s.debugFunc = "chantimer"
			t.Cb()
		}
	}
}

func (s *Skeleton) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	if s.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.AfterFunc(d, cb)
}

func (s *Skeleton) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	if s.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.CronFunc(cronExpr, cb)
}

func (s *Skeleton) Go(f func(), cb func()) {
	if s.GoLen == 0 {
		panic("invalid GoLen")
	}

	s.g.Go(f, cb)
}

func (s *Skeleton) NewLinearContext() *g.LinearContext {
	if s.GoLen == 0 {
		panic("invalid GoLen")
	}

	return s.g.NewLinearContext()
}

func (s *Skeleton) AsynCall(server *chanrpc.Server, id interface{}, args ...interface{}) {
	if s.AsynCallLen == 0 {
		panic("invalid AsynCallLen")
	}

	s.client.Attach(server)
	s.client.AsynCall(id, args...)
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.ChanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.server.Register(id, f)
}

func (s *Skeleton) RegisterCommand(name string, help string, f interface{}) {
	console.Register(name, help, f, s.commandServer)
}

func (s *Skeleton) ThreadRun() {
	//s.AfterFunc(time.Microsecond*30, s.ThreadRun)
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}
	}()
	now := time.Now().UnixNano() / 1e6
	diff := now - s.strartTime - s.runTime
	s.runTime += diff
	if diff < 0 {
		s.strartTime = now
		diff = 0
	}
	if diff < 100000 { //大于100秒，本次循环跳过
		if s.Engine != nil {
			s.Engine.Update(diff)
		}
	} else {
		log.Error("%v ThreadRun diff %v", s.Name, diff)
	}
}
func (s *Skeleton) Update(diff int64) {

}

func (s *Skeleton) Debug() {
	log.Debug("%v: %v>>%v", s.Name, s.debugFunc, s.debugData)
}

func NewSkeleton() *Skeleton {
	skeleton := &Skeleton{
		GoLen:              5000,
		TimerDispatcherLen: 5000,
		AsynCallLen:        50000,
		ChanRPCServer:      chanrpc.NewServer(50000),
	}
	skeleton.Init()
	return skeleton
}
