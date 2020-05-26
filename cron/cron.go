package cron

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// 定时任务计划
/*
- spec，传入 cron 时间设置
- job，对应执行的任务
- flag，管道用于关闭任务
*/
func StartJob(spec string, job cron.Job, shut chan int) {
	logger := &CLog{clog: log.New()}
	logger.clog.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(logger)))

	c.AddJob(spec, job)

	// 启动执行任务
	c.Start()
	// 退出时关闭计划任务
	defer c.Stop()

	select {
	case <-shut:
		return
	}
}

func StopJob(shut chan int) {
	shut <- 0
}

type CLog struct {
	clog *log.Logger
}

func (l *CLog) Info(msg string, keysAndValues ...interface{}) {
	l.clog.WithFields(log.Fields{
		"data": keysAndValues,
	}).Info(msg)
}

func (l *CLog) Error(err error, msg string, keysAndValues ...interface{}) {
	l.clog.WithFields(log.Fields{
		"msg":  msg,
		"data": keysAndValues,
	}).Warn(msg)
}
