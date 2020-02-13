package gate

import (
	"encoding/binary"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/module"
	"github.com/wudiliujie/common/network"
	"github.com/wudiliujie/common/pool"

	"net"
	"reflect"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

var sessionId int32 = 0

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	// agent
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	ChanRPCLen         int
	NetEvent           *NetEvent
}
type NetEvent struct {
	OnAgentInit    func(Agent)
	OnAgentDestroy func(Agent)
	OnReceiveMsg   func(agent Agent, data []byte)
	Processor      network.IProcessor
}

func newAgent(conn network.Conn, g *NetEvent, tag int64) network.Agent {
	a := &agent{conn: conn}
	a.tag = tag
	a.netEvent = g
	//如果是json，消息类型修改为：TextMessage     websocket.BinaryMessage
	a.id = atomic.AddInt32(&sessionId, 1)
	if a.netEvent.OnAgentInit != nil {
		a.netEvent.OnAgentInit(a)
	}
	return a
}

func (gate *Gate) GetName() string {
	return "gate"
}
func (gate *Gate) Run(closeSig chan bool) {

	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			return newAgent(conn, gate.NetEvent, int64(0))
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			return newAgent(conn, gate.NetEvent, int64(0))
		}
	}

	if wsServer != nil {
		log.Debug("start ws")
		wsServer.Start()
	}
	if tcpServer != nil {
		log.Debug("start tcp")
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}
func (gate *Gate) Debug()     {}

type agent struct {
	conn       network.Conn
	netEvent   *NetEvent
	userData   interface{}
	module     module.Module
	id         int32
	isLogin    bool
	heartTimes int32
	agentType  int32
	tag        int64
	auto       bool
	idx        uint16 //需要删除缓存的包
}

func (a *agent) GetIdx() uint16 {
	return a.idx
}

func (a *agent) GetAutoReconnect() bool {
	return a.auto
}

func (a *agent) SetAutoReconnect(v bool) {
	a.auto = v
}

func (a *agent) SetTag(tag int64) {
	a.tag = tag
}

func (a *agent) GetTag() int64 {
	return a.tag
}

func (a *agent) ResetHeart() {
	a.heartTimes = 0
}

func (a *agent) AddHeart() {
	a.heartTimes++
}

func (a *agent) GetHeart() int32 {
	return a.heartTimes
}

func (a *agent) Call1(userId int64, pck interface{}) (string, error) {
	panic("implement me")
}

func (a *agent) Run() {
	closeSig := make(chan bool, 1)
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}

		closeSig <- true
	}()

	for {
		data, err := a.conn.ReadMsg()

		//a := string(data)
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}
		if a.netEvent.OnReceiveMsg != nil {
			a.netEvent.OnReceiveMsg(a, data)
		}
	}
}

func (a *agent) OnClose() {
	if a.netEvent.OnAgentDestroy != nil {
		//log.Debug("close agent")
		a.netEvent.OnAgentDestroy(a)
	}
}

func (a *agent) WriteMsg(idx uint16, msg network.IMessage) {
	//这里缓存
	buff := pool.GetBytesLen(2)
	binary.BigEndian.PutUint16(buff, idx)
	var err error
	if a.netEvent.Processor != nil {
		buff, err = a.netEvent.Processor.MarshalBytes(buff, msg)
		if err != nil {
			log.Error("marshal message %v error: %v  %s", reflect.TypeOf(msg), err, debug.Stack())
			return
		}
	}
	err = a.conn.WriteMsg(buff)
	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}

}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}
func (a *agent) Ip() string {
	aa := strings.Split(a.conn.RemoteAddr().String(), ":")
	if len(aa) >= 1 {
		return aa[0]
	}
	return "0.0.0.0"
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) GetId() int32 {
	return a.id
}

func (a *agent) SetLogin() {
	a.isLogin = true
}
func (a *agent) GetLogin() bool {
	return a.isLogin
}
