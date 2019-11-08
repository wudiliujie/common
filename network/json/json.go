package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wudiliujie/common/chanrpc"
	"github.com/wudiliujie/common/convert"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/network"
	"reflect"
)

type Processor struct {
	msgInfo map[uint16]*MsgInfo
	msgID   map[reflect.Type]uint16
	MsgType network.MsgType
}

func (p *Processor) GetMsgType() network.MsgType {
	return p.MsgType
}

type MsgInfo struct {
	msgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      string
	msgRawData json.RawMessage
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.msgInfo = make(map[uint16]*MsgInfo)
	p.msgID = make(map[reflect.Type]uint16)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(id uint16, msg interface{}) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	//msgID := msgType.Elem().Name()
	//if msgID == "" {
	//	log.Fatal("unnamed json message")
	//}
	if _, ok := p.msgInfo[id]; ok {
		log.Fatal("message %v is already registered", id)
	}
	p.msgID[msgType] = id
	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[id] = i
	return id
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg interface{}, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := uint16(p.GetMsgId(msg))
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRouter = msgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msgId uint16, msgHandler MsgHandler) {

	i, ok := p.msgInfo[msgId]
	if !ok {
		log.Fatal("message %v not registered", msgId)
	}
	i.msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(msgID string, msgRawHandler MsgHandler) {
	//i, ok := p.msgInfo[msgID]
	//if !ok {
	//	log.Fatal("message %v not registered", msgID)
	//}
	//
	//i.msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {

	// json
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgID := uint16(p.GetMsgId(msg))
	i, ok := p.msgInfo[msgID]
	if !ok {
		return fmt.Errorf("message %v not registered", msgID)
	}
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
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if len(m) != 2 {
		return nil, errors.New("invalid json data")
	}
	i, ok := p.msgInfo[uint16(convert.ToInt32(string(m["c"])))]
	if !ok {
		return nil, fmt.Errorf("message %v not registered", m["c"])
	} // msg
	msg := reflect.New(i.msgType.Elem()).Interface()
	return msg, json.Unmarshal(m["b"], msg)
	panic("bug")
}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return nil, errors.New("json message pointer required")
	}
	msgID := uint16(p.GetMsgId(msg))
	if _, ok := p.msgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}
	m := map[string]interface{}{"b": msg, "c": msgID}
	data, err := json.Marshal(m)
	return [][]byte{data}, err
}
func (p *Processor) GetMsgId(pck interface{}) int32 {
	msgType := reflect.TypeOf(pck)
	id, ok := p.msgID[msgType]
	if !ok {
		return 0
	}
	return int32(id)
}
