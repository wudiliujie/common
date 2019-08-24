package kafka_producer

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/satori/go.uuid"
	"github.com/wudiliujie/common/log"
	"strings"
	"time"
)

type Kafka_Producer struct {
	pProducer    sarama.AsyncProducer
	kafka_iplist string
	topics       string
}
type Kafka_Producer_interface interface {
	AsyncProducer(vjson string)
	Close()
	Kafka_insertDeviceRamLog(data map[interface{}]interface{}, Operation, Table string)
}

func (self *Kafka_Producer) Close() {
	self.pProducer.Close()
}

func New_Kafka(iplist, topic string) (*Kafka_Producer, error) {
	outres := &Kafka_Producer{
		kafka_iplist: iplist,
		topics:       topic,
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewAsyncProducer(strings.Split(outres.kafka_iplist, ","), config)
	if err != nil {
		return nil, err
	}

	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					_, err := err.Msg.Value.Encode()
					if err != nil {
						log.Error("AsyncProducer   %s", err)
					} else {

					}
				}
			case <-success:
			}
		}
	}(p)
	outres.pProducer = p

	return outres, nil
}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func (self *Kafka_Producer) AsyncProducer(vjson string, topic string) {
	//now := time.Now()
	if topic == "" {
		topic = self.topics
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(vjson),
	}
	self.pProducer.Input() <- msg
	//fmt.Println("...................", time.Now().UnixNano()/1e6-now.UnixNano()/1e6)
}

//获取当前时间戳
func GetLocalTimeMircTime() int64 {
	return time.Now().UTC().UnixNano()/1000000 + 8*3600*1000
}

type Kafka_SJsonDeviceRamLog struct {
	Operation string      `json:"operation"`
	Sql       interface{} `json:"sql"`
	Table     string      `json:"table"`
	SrcHost   uint32      `json:"src_host"`
	Id        string      `json:"id"`
	GameId    int32       `json:"game_id"`
}

func GetUUID() string {
	str_uuid := uuid.NewV4()
	return str_uuid.String()
}
func (self *Kafka_Producer) Kafka_InstertLog(data map[string]interface{}, Operation, Table string) (id string) {
	str_uuid := GetUUID()
	self.Kafka_InstertLogById(str_uuid, data, Operation, Table)
	return str_uuid
}
func (self *Kafka_Producer) Kafka_InstertLogById(id string, data map[string]interface{}, Operation, Table string) {
	tmp := &Kafka_SJsonDeviceRamLog{
		Operation: Operation,
		Table:     Table,
		SrcHost:   1,
	}
	data["id"] = id
	data["create_time"] = GetLocalTimeMircTime()
	//
	tmp.Sql = data
	data1, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println("kfk log插入失败", string(data1))
	} else {
		//fmt.Println("kfk log ",self.topics+"_"+Table, string(data1))
		self.AsyncProducer(string(data1), self.topics+"_"+Table)
	}
}
func (self *Kafka_Producer) Kafka_InstertLogArry(data []map[string]interface{}, Operation, Table string) {
	tmp := &Kafka_SJsonDeviceRamLog{
		Operation: Operation,
		Table:     Table,
		SrcHost:   1,
	}
	//
	tmp.Sql = data
	tmp.Id = GetUUID()

	data1, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println("kfk log插入失败", string(data1))
	} else {
		//fmt.Println("kfk log ",self.topics+"_"+Table, string(data1))
		self.AsyncProducer(string(data1), self.topics+"_"+Table)
	}
}
