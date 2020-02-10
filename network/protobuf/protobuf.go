package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/network"
	"github.com/wudiliujie/common/pool"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	msgInfo      map[uint16]func() network.IMessage
	msgType      network.MsgType
}

func (p *Processor) GetMsgType() network.MsgType {
	return p.msgType
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.msgType = network.MsgType_Proto
	p.msgInfo = make(map[uint16]func() network.IMessage)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(id uint16, f func() network.IMessage) {
	if _, ok := p.msgInfo[id]; ok {
		log.Fatal("message %s is already registered", id)
	}
	p.msgInfo[id] = f

}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (network.IMessage, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	// id
	var id uint16
	if p.littleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}

	// msg
	i := p.msgInfo[id]
	if i != nil {
		msg := i()
		return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
	} else {
		return nil, errors.New(fmt.Sprintf("id:%v   不存在", id))

	}

}

// goroutine safe
func (p *Processor) Marshal(msg network.IMessage) ([]byte, error) {

	id := pool.GetBytesLen(2)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, msg.GetId())
	} else {
		binary.BigEndian.PutUint16(id, msg.GetId())
	}
	//这里可用缓存？
	buff := proto.NewBuffer(id)
	err := buff.Marshal(msg.(proto.Message))
	return buff.Bytes(), err
}
