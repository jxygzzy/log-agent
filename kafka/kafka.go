package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client  sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka:producer closed,err:%v", err)
		return err
	}
	Client = client
	MsgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	return nil
}

func sendMsg() {
	for msg := range MsgChan {
		pid, offset, err := Client.SendMessage(msg)
		if err != nil {
			logrus.Warn("send msg failed,err:%v", err)
		}
		logrus.Infof("send msg to kafka success,pid:%v offset:%v", pid, offset)
	}
}
