package main

import (
	"strings"
	"time"

	"logagent/kafka"
	"logagent/tailfile"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int    `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

func run() error {
	for {
		line, ok := <-tailfile.TailTask.Lines
		if !ok {
			logrus.Warn("tail file close reopen, filename:%s\n", tailfile.TailTask.Filename)
			time.Sleep(time.Second)
			continue
		}
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.MsgChan <- msg
	}
}

// 日志收集客户端
// 收集指定目录下的日志 -> 发送至kafka
// 使用tailf读取日志文件

func main() {
	var config = new(Config)
	err := ini.MapTo(config, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed,err:%v", err)
	}

	err = kafka.Init(strings.Split(config.KafkaConfig.Address, ","), config.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init Kafka success!")

	err = tailfile.Init(config.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Error("init tail failed,err:%v", err)
		return
	}
	logrus.Info("init tail success!")

	err = run()
	if err != nil {
		logrus.Error("run failed,err:%v", err)
		return
	}
}
