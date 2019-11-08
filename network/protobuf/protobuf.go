package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wudiliujie/common/chanrpc"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/network"
	"github.com/wudiliujie/common/pool"
	"math"
	"reflect"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	msgInfo      map[uint16]*MsgInfo
	msgID        map[reflect.Type]uint16
	msgType      network.MsgType
}

func (p *Processor) GetMsgType() network.MsgType {
	return p.msgType
}

type MsgInfo struct {
	msgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      uint16
	msgRawData []byte
}

func NewProcessor(maxid int) *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.msgType = network.MsgType_Proto
	p.msgID = make(map[reflect.Type]uint16)
	p.msgInfo = make(map[uint16]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(id uint16, msg proto.Message) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if _, ok := p.msgID[msgType]; ok {
		log.Fatal("message %s is already registered", msgType)
	}
	if len(p.msgInfo) >= math.MaxUint16 {
		log.Fatal("too many protobuf messages (max = %v)", math.MaxUint16)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[id] = i
	//id := uint16(len(p.msgInfo) - 1)
	p.msgID[msgType] = id
	return id
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatal("message %s not registered", msgType)
	}

	p.msgInfo[id].msgRouter = msgRouter
}

//// It's dangerous to call the method on routing or marshaling (unmarshaling)
//func (p *Processor) SetHandler(msg proto.Message, msgHandler MsgHandler) {
//	msgType := reflect.TypeOf(msg)
//	id, ok := p.msgID[msgType]
//	if !ok {
//		log.Fatal("message %s not registered", msgType)
//	}
//
//	p.msgInfo[id].msgHandler = msgHandler
//}
// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msgId uint16, msgHandler MsgHandler) {

	i, ok := p.msgInfo[msgId]
	if !ok {
		log.Fatal("message %v not registered", msgId)
	}
	i.msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(id uint16, msgRawHandler MsgHandler) {
	if id >= uint16(len(p.msgInfo)) {
		log.Fatal("message id %v not registered", id)
	}

	p.msgInfo[id].msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	// raw
	if msgRaw, ok := msg.(MsgRaw); ok {
		if msgRaw.msgID >= uint16(len(p.msgInfo)) {
			return fmt.Errorf("message id %v not registered", msgRaw.msgID)
		}
		i := p.msgInfo[msgRaw.msgID]
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}

	// protobuf
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		return fmt.Errorf("message %s not registered", msgType)
	}
	i := p.msgInfo[id]
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgType, msg, userData)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
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
	//if id >= uint16(len(p.msgInfo)) {
	//	return nil, fmt.Errorf("message id %v not registered", id)
	//}

	// msg
	i := p.msgInfo[id]
	if i != nil {
		if i.msgRawHandler != nil {
			return MsgRaw{id, data[2:]}, nil
		} else {
			msg := reflect.New(i.msgType.Elem()).Interface()
			return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("id:%v   不存在", id))

	}

}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([]byte, error) {
	msgType := reflect.TypeOf(msg)

	// id
	_id, ok := p.msgID[msgType]
	if !ok {
		err := fmt.Errorf("message %s not registered", msgType)
		return nil, err
	}

	id := pool.GetBytesLen(2)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, _id)
	} else {
		binary.BigEndian.PutUint16(id, _id)
	}
	//这里可用缓存？
	buff := proto.NewBuffer(id)
	err := buff.Marshal(msg.(proto.Message))
	return buff.Bytes(), err
}

// goroutine safe
func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
	for id, i := range p.msgInfo {
		f(uint16(id), i.msgType)
	}
}
func (p *Processor) GetMsgId(pck interface{}) int32 {
	msgType := reflect.TypeOf(pck)
	id, ok := p.msgID[msgType]
	if !ok {
		return 0
	}
	return int32(id)
}
