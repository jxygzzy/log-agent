package tailfile

import (
	"context"
	"logagent/kafka"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

type tailTask struct {
	path   string
	topic  string
	obj    *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func (t *tailTask) Init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.obj, err = tail.TailFile(t.path, config)
	return
}

func (t *tailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			logrus.Warnf("the task for path:%s is stop...", t.path)
			t.obj.Cleanup()
			return
		case line, ok := <-t.obj.Lines:
			if !ok {
				logrus.Warn("tail file close reopen, filename:%s\n", t.path)
				time.Sleep(time.Second)
				continue
			}
			if len(strings.Trim(line.Text, "\r")) == 0 {
				// 空行跳过
				continue
			}
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.ToMsgChan(msg)
		}

	}
}
