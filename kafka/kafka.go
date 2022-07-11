package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka:producer closed, err:%v", err)
		return
	}
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	return
}

func sendMsg() {
	for msg := range msgChan {
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			logrus.Warn("send msg failed, err:%v", err)
		}
		logrus.Infof("send msg to kafka success, pid:%v offset:%v", pid, offset)
	}
}

func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
