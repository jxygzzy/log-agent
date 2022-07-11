package tailfile

import (
	"context"
	"logagent/common"

	"github.com/sirupsen/logrus"
)

var (
	ttManager *tailTaskManager
)

type tailTaskManager struct {
	tailTaskMap     map[string]*tailTask
	collectEntyList []common.CollectEnty
	confChan        chan []common.CollectEnty
}

func Init(allConf []common.CollectEnty) error {
	ttManager = &tailTaskManager{
		tailTaskMap:     make(map[string]*tailTask, 16),
		collectEntyList: allConf,
		confChan:        make(chan []common.CollectEnty),
	}
	for _, conf := range allConf {
		ctx, cancel := context.WithCancel(context.Background())
		tt := tailTask{
			path:   conf.Path,
			topic:  conf.Topic,
			ctx:    ctx,
			cancel: cancel,
		}
		err := tt.Init()
		if err != nil {
			logrus.Error("tailfile: create tailTask for path:%s failed,err:%v", tt.path, err)
			continue
		}
		ttManager.tailTaskMap[tt.topic+tt.path] = &tt
		go tt.run()
	}
	go ttManager.watch()
	return nil
}

func (tm *tailTaskManager) watch() {
	for newConf := range tm.confChan {
		newMap := make(map[string]bool)
		for _, conf := range newConf {
			newMap[conf.Topic+conf.Path] = true
			if tm.isExist(conf) {
				continue
			}
			// 新增任务
			tt := tailTask{
				path:  conf.Path,
				topic: conf.Topic,
			}
			err := tt.Init()
			if err != nil {
				logrus.Error("tailfile: create tailTask for path:%s failed,err:%v", tt.path, err)
				continue
			}
			ttManager.tailTaskMap[tt.topic+tt.path] = &tt
			go tt.run()
		}
		for key, task := range tm.tailTaskMap {
			if _, ok := newMap[key]; !ok {
				// 删除
				task.cancel()
				delete(tm.tailTaskMap, key)
			}
		}
	}
}
func (tm *tailTaskManager) isExist(conf common.CollectEnty) bool {
	if _, ok := tm.tailTaskMap[conf.Topic+conf.Path]; ok {
		return true
	}
	return false
}

func SetNewConf(conf []common.CollectEnty) {
	ttManager.confChan <- conf
}
