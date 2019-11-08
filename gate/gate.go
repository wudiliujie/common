package gate

import (
	"encoding/binary"
	"errors"
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
	Processor       network.Processor

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
	OnAgentInit        func(Agent)
	OnAgentDestroy     func(Agent)
	OnReceiveMsg       func(agent Agent, msgId int32, pck interface{})
}

func newAgent(conn network.Conn, tag interface{}) network.Agent {
	a := &agent{conn: conn}
	a.tag = tag
	a.gate = ThisGate
	if ThisGate.Processor.GetMsgType() == network.MsgType_Json {
		conn.LocalAddr()
	}
	//如果是json，消息类型修改为：TextMessage     websocket.BinaryMessage
	a.id = atomic.AddInt32(&sessionId, 1)
	if ThisGate == nil {

	}
	if ThisGate.OnAgentInit != nil {
		ThisGate.OnAgentInit(a)
	}
	return a
}

var ThisGate *Gate

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
			return newAgent(conn, int64(0))
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
			return newAgent(conn, int64(0))
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
	gate       *Gate
	userData   interface{}
	module     module.Module
	id         int32
	isLogin    bool
	heartTimes int32
	agentType  int32
	tag        interface{}
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

func (a *agent) SetTag(tag interface{}) {
	a.tag = tag
}

func (a *agent) GetTag() interface{} {
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

	handleMsgData := func(args []interface{}) error {
		if a.gate.Processor != nil {
			data := args[0].([]byte)
			if len(data) < 4 {
				log.Debug("长度不够")
				return errors.New("长度不够")
			}
			//log.Debug("*******************%v", string(data))
			//索引:%v
			a.idx = binary.BigEndian.Uint16(data)
			msg, err := a.gate.Processor.Unmarshal(data[2:])
			//msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				return err
			}
			msgId := a.gate.Processor.GetMsgId(msg)
			if a.gate.OnReceiveMsg != nil {
				a.gate.OnReceiveMsg(a, msgId, msg)
			}
		}
		return nil
	}
	for {
		data, err := a.conn.ReadMsg()

		//a := string(data)
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		err = handleMsgData([]interface{}{data})
		if err != nil {
			log.Debug("handle message: %v", err)
			//break
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.OnAgentDestroy != nil {
		//log.Debug("close agent")
		a.gate.OnAgentDestroy(a)
	}
}

func (a *agent) WriteMsg(idx uint16, msg interface{}) {
	//这里缓存
	id := pool.GetBytesLen(2)
	binary.BigEndian.PutUint16(id, idx)
	if a.gate.Processor != nil {
		switch msg.(type) {
		case string:
			id = append(id, []byte(msg.(string))...)
			break
		case []uint8:
			id = append(id, msg.([]uint8)...)
			break
		default:
			data, err := a.gate.Processor.Marshal(msg)
			if err != nil {
				log.Error("marshal message %v error: %v  %s", reflect.TypeOf(msg), err, debug.Stack())
				return
			}

			id = append(id, data...)
			pool.PutBytes(data)
		}

	}
	err := a.conn.WriteMsg(id)
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
func (a *agent) Gate() *Gate {
	return a.gate
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
