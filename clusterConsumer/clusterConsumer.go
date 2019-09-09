package clusterConsumer

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/module"
	"strings"
)

const ConsumerMsgEvent = "ConsumerMsgEvent"

type ClusterConsumer struct {
	brokers  []string
	topics   []string
	groupId  string
	consumer *cluster.Consumer
	closeSig chan bool
}

func NewClusterConsumer(brokers string, topics string, groupId string) *ClusterConsumer {
	c := &ClusterConsumer{}
	c.brokers = strings.Split(brokers, ",")
	c.topics = strings.Split(topics, ",")
	c.groupId = groupId
	c.closeSig = make(chan bool, 1)
	return c
}
func (c *ClusterConsumer) Init() error {
	var err error
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	c.consumer, err = cluster.NewConsumer(c.brokers, c.groupId, c.topics, config)
	return err
}
func (c *ClusterConsumer) Run() {
	for {
		select {
		case <-c.closeSig:
			log.Debug("close")
			return
		case msg, ok := <-c.consumer.Messages():
			if ok {
				log.Debug("msg:%v %v", msg.Offset, msg.Topic)
				module.OnChanRpcEvent(ConsumerMsgEvent, c, msg)
				//consumer.MarkOffset(msg, "")
			}
			break
		case ntf := <-c.consumer.Notifications():
			if ntf != nil {
				log.Debug("notifications %v", ntf)
			}
			break
		case err := <-c.consumer.Errors():
			if err != nil {
				log.Error("kfk进程异常 %v", err)
			}
			break
		}
	}
}
func (c *ClusterConsumer) MarkOffset(msg *sarama.ConsumerMessage) {
	c.consumer.MarkOffset(msg, "")
}
func (c *ClusterConsumer) Close() {
	c.closeSig <- true
}
