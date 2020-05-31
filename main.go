package main

import (
	"flag"
	"github.com/PointStoneTeam/cugsync/manager"
	"github.com/PointStoneTeam/cugsync/routers"
	"github.com/PointStoneTeam/cugsync/rsync"
	"github.com/PointStoneTeam/cugsync/setting"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// @title CUGSYNC
// @version 0.0.0
// @description 中国地质大学（武汉）开源镜像站

// @contact.name PointStone
// @contact.url https://mirrors.cug.edu.cn
// @contact.qq 1011615604

// @host mirrors.cug.edu.cn/api/
func main() {

	confPath := flag.String("conf", "conf/config.json", "指定配置文件路径")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 加载用户配置
	if err := setting.LoadUserConfig(*confPath); err != nil {
		log.Fatal(err)
	}
	// 判断 rsync 是否安装
	if hasRsync, info := rsync.CheckRsync(); !hasRsync {
		log.Fatal(info)
	}
	// 导入默认配置的 Jobs，默认启动从配置文件读到的任务
	if jobList, err := manager.GetDefaultJob(setting.GetDefaultJobPath()); err != nil {
		log.WithField("err", err).Warnln("没有任务计划文件，启动后将无初始任务")
	} else {
		manager.InitJobs(jobList)
	}

	// 处理端口绑定
	Addr := setting.GetBindAddr(setting.SyncSet.Server.BindGlobal, setting.SyncSet.Server.Port)

	// 启动服务器
	server := &http.Server{
		Addr:           Addr,
		Handler:        routers.InitRouter(),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	panic(server.ListenAndServe())
}
