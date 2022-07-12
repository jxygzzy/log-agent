package etcd

import (
	"context"
	"encoding/json"
	"logagent/common"
	"logagent/tailfile"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/storage/storagepb"
	"github.com/sirupsen/logrus"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		logrus.Error("connect to etcd failed, err:%v", err)
		return
	}
	return
}

func GetConfig(key string) (collectEntyList []common.CollectEnty, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s failed, err:%v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warnf("get len:0 conf from etcd by key:%s", key)
		return
	}
	ret := resp.Kvs[0].Value
	err = json.Unmarshal(ret, &collectEntyList)
	if err != nil {
		logrus.Errorf("json unmarshal failed, err:%v", err)
		return
	}
	return
}

func WatchConf(key string) {
	for {
		watchChan := client.Watch(context.Background(), key)
		for resp := range watchChan {
			for _, evt := range resp.Events {
				var newConf []common.CollectEnty
				if evt.Type == storagepb.DELETE {
					tailfile.SetNewConf(newConf)
					continue
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal new conf failed, err:%v", err)
					continue
				}
				tailfile.SetNewConf(newConf)
			}
		}
	}
}
