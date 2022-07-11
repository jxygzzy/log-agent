package etcd

import (
	"testing"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
)

func TestEtcdPut(t *testing.T) {
	err := Init([]string{"127.0.0.1:2379"})
	if err != nil {
		t.Error(err)
	}
	client.Put(context.Background(), "collect_log_conf", `[{"path":"c:/logs/xx.log","topic":"web_log"},{"path":"c:/logs/xx1.log","topic":"shopping_log"}]`)
}
