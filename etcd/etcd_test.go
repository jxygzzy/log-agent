package etcd

import (
	"fmt"
	"testing"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
)

func TestEtcdPut(t *testing.T) {
	err := Init([]string{"127.0.0.1:2379"})
	if err != nil {
		t.Error(err)
	}
	client.Put(context.Background(), "collect_log_192.168.1.36_conf", `[{"path":"c:/logs/xx.log","topic":"web_log"},{"path":"c:/logs/xx1.log","topic":"shopping_log"}]`)
	// client.Delete(context.Background(), "collect_log_192.168.1.36_conf")
	resp, _ := client.Get(context.Background(), "collect_log_192.168.1.36_conf")
	fmt.Printf("resp:%s", resp.Kvs[0].Value)
}
