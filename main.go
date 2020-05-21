package main

import (
	"flag"
	"github.com/PointStoneTeam/cugsync/conf"
	"github.com/PointStoneTeam/cugsync/routers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {

	confPath := flag.String("conf", "config.json", "指定配置文件路径")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 加载用户配置
	if err := conf.LoadUserConfig(*confPath); err != nil {
		log.Fatal(err)
	}

	// 处理端口绑定
	Addr := conf.GetBindAddr(conf.SyncSet.Server.BindGlobal, conf.SyncSet.Server.Port)

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
