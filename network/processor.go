package network

type MsgType int32

const (
	MsgType_Json MsgType = iota
	MsgType_Proto
)

type IMessage interface {
	GetId() uint16
}
type IProcessor interface {

	// must goroutine safe
	Unmarshal(data []byte) (IMessage, error)
	// must goroutine safe
	Marshal(msg IMessage) ([]byte, error)
	GetMsgType() MsgType
	MarshalBytes(buffer []byte, msg IMessage) ([]byte, error)
}
