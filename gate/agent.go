package gate

import (
	"github.com/wudiliujie/common/network"
	"net"
)

type Agent interface {
	network.Agent
	GetId() int32
	WriteMsg(idx uint16, msg network.IMessage)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Ip() string
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})

	GetLogin() bool
	SetLogin()
	Call1(userId int64, pck interface{}) (string, error)
	AddHeart()
	ResetHeart()
	GetHeart() int32
	GetTag() int64
	SetTag(int64)
	GetIdx() uint16
}
