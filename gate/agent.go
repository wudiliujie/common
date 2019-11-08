package gate

import (
	"github.com/wudiliujie/common/network"
	"net"
)

type Agent interface {
	network.Agent
	GetId() int32
	WriteMsg(idx uint16, msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Ip() string
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	Gate() *Gate
	GetLogin() bool
	SetLogin()
	Call1(userId int64, pck interface{}) (string, error)
	AddHeart()
	ResetHeart()
	GetHeart() int32
	GetTag() interface{}
	SetTag(interface{})
	GetIdx() uint16
}
