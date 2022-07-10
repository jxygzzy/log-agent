package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var TailTask *tail.Tail

func Init(filename string) error {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tailTask, err := tail.TailFile(filename, config)
	if err != nil {
		logrus.Error("tailfile: create tailTask for path:%s failed,err:%v", err)
		return err
	}
	TailTask = tailTask
	return nil
}
