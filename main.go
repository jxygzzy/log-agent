package main

import (
	"fmt"
	"strings"

	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	KafkaConfig `ini:"kafka"`
	EtcdConfig  `ini:"etcd"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

// 日志收集客户端
// 收集指定目录下的日志 -> 发送至kafka
// 使用tailf读取日志文件

func main() {
	ip, err := common.GetLocalIp()
	if err != nil {
		logrus.Error("get local ip failed, err:%v", err)
		return
	}
	logrus.Infof("get local ip success, ip:%s", ip)
	var config = new(Config)
	err = ini.MapTo(config, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err:%v", err)
		return
	}
	err = kafka.Init(strings.Split(config.KafkaConfig.Address, ","), config.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init Kafka success!")

	err = etcd.Init(strings.Split(config.EtcdConfig.Address, ","))
	if err != nil {
		logrus.Error("init etcd failed, err:%v", err)
		return
	}
	logrus.Info("init etcd success!")

	collectKey := fmt.Sprintf(config.EtcdConfig.CollectKey, ip)
	allConf, err := etcd.GetConfig(collectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err:%v", err)
	}
	logrus.Info("get conf from etcd success!")
	go etcd.WatchConf(collectKey)
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Error("init tail failed, err:%v", err)
		return
	}
	logrus.Info("init tail success!")
	run()
}

func run() {
	select {}
}
