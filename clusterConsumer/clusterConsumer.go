package clusterConsumer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/wudiliujie/common/log"
	"github.com/wudiliujie/common/module"
	"strings"
	"sync"
)

const ConsumerMsgEvent = "ConsumerMsgEvent"

type ClusterConsumer struct {
	brokers           []string
	topics            []string
	groupId           string
	consumer          sarama.ConsumerGroup
	closeSig          chan bool
	partitionConsumer []sarama.PartitionConsumer
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
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V1_0_0_0
	c.consumer, err = sarama.NewConsumerGroup(c.brokers, c.groupId+"aaa", config)
	return err
}
func (c *ClusterConsumer) Run() {
	consumer := Consumer{
		ready: make(chan bool),
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumer.Consume(ctx, c.topics, &consumer); err != nil {
				log.Fatal("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready // Await till the consumer has been set up
	select {
	case <-ctx.Done():
		log.Debug("terminating: context cancelled")
	case <-c.closeSig:
		log.Debug("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := c.consumer.Close(); err != nil {
		log.Fatal("Error closing client: %v", err)
	}
}

func (c *ClusterConsumer) Close() {
	c.closeSig <- true
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		module.OnChanRpcEvent(ConsumerMsgEvent, message)
		session.MarkMessage(message, "")
	}

	return nil
}
