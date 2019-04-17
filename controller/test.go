package controller

import (
	"fmt"
	"redisdemo/model"
	"sync"
	"time"

	"redisdemo/engine"
	"redisdemo/tools"
)

type TestOb struct {
	MaxWorker int
	TestChan  chan string
	RedisKey  string
	IsClose   bool
	mu        sync.RWMutex
}

func (t *TestOb) Test() {

	go func() {
		for {
			select {
			case d := <-t.TestChan:
				tools.Log.Info("当前channel 数据长度为：", len(t.TestChan))
				time.Sleep(5 * time.Second) //模拟处理任务耗时
				fmt.Println(d, "-----------")
				test := new(model.Test)
				test.Name = d
				if err := test.InsertOne(); err != nil {
					tools.Log.Error(err)
				}
			}
		}
	}()
	cli := engine.RedisCli
	defer cli.Close()
FORLABEL:
	for {
		select {
		case <-time.Tick(5 * time.Second): //五秒检查一次队列中是否有任务
		AGAIN:
			t.mu.Lock()
			if t.IsClose {
				t.mu.Unlock()
				break FORLABEL
			}
			t.mu.Unlock()
			count, err := cli.LLen(t.RedisKey).Result()
			if err != nil {
				tools.Log.Warn("get llen error ", err.Error())
				continue
			}
			if count < 1 {
				continue
			}

			var i int64
			for i = 0; i < count; i++ {
				data, err := cli.LPop(t.RedisKey).Result()
				if err != nil {
					tools.Log.Warn("get lpop error ", err.Error())
					continue
				}

				t.TestChan <- data
			}
			goto AGAIN
		}

	}

}

func (t *TestOb) StopTest() {
	t.mu.Lock()
	t.IsClose = true
	t.mu.Unlock()
LABEL:
	for {
		select {
		case <-time.Tick(1 * time.Second):
			if len(t.TestChan) < 1 {
				break LABEL
			}
		}
	}
}
