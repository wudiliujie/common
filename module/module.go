package module

import (
	"github.com/wudiliujie/common/chanrpc"
	"github.com/wudiliujie/common/log"
	"runtime"
	"sync"
)

var LenStackBuf = 4096

type Module interface {
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
	Debug()
}

type chanRpcEvent struct {
	server *chanrpc.Server
	f      interface{}
}

type module struct {
	mi       Module
	closeSig chan bool
	wg       sync.WaitGroup
}

var mods []*module
var mapEvent = make(map[string][]*chanRpcEvent)

func Register(mi Module) {
	m := new(module)
	m.mi = mi
	m.closeSig = make(chan bool, 1)

	mods = append(mods, m)
}
func GetAllModule() []Module {
	ret := make([]Module, 0)
	for _, v := range mods {
		ret = append(ret, v.mi)
	}
	return ret
}

func Init() {
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit()
	}

	for i := 0; i < len(mods); i++ {
		m := mods[i]
		m.wg.Add(1)
		go run(m)
	}
}

func Destroy() {
	for i := len(mods) - 1; i >= 0; i-- {
		m := mods[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func run(m *module) {
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *module) {
	defer func() {
		if r := recover(); r != nil {
			if LenStackBuf > 0 {
				buf := make([]byte, LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("destroy %v", r)
			}
		}
	}()
	log.Debug("%v", m)
	m.mi.OnDestroy()
}

func RegisterChanRpcEvent(event string, _server *chanrpc.Server, _f interface{}) {
	_server.Register(event, _f)
	mapEvent[event] = append(mapEvent[event], &chanRpcEvent{server: _server, f: _f})
}
func OnChanRpcEvent(event string, args ...interface{}) {
	v, ok := mapEvent[event]
	if ok {
		for _, f := range v {
			f.server.Go(event, args...)
		}
	}
}
