package main

import (
	"os"
	"os/signal"
	"syscall"

	"redisdemo/controller"
	"redisdemo/engine"
	"redisdemo/tools"
)

func main() {

	tools.InitLog()
	tools.InitConfig()

	engine.InitRedisEngine()
	t := new(controller.TestOb)
	t.RedisKey = "test_queue"
	t.MaxWorker = 10
	go t.Test()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	tools.Log.Warn(<-ch, " Signal received. Shutting down...")
	t.StopTest()
}
