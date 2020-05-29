package setting

import (
	"encoding/json"
	"fmt"
	"github.com/PointStoneTeam/cugsync/pkg/file"
	log "github.com/sirupsen/logrus"
	"time"
)

// 服务器设置
type Server struct {
	Port         int `json:"port"`
	RefreshTime  int `json:"refresh_time"` //单位为分钟
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BindGlobal   bool   `json:"bind_global"`
	DistPATH     string `json:"dist_path"` // 静态文件目录
}

var defaultServerSetting = &Server{
	Port:         8000,
	RefreshTime:  10,
	ReadTimeout:  60,
	WriteTimeout: 60,
	BindGlobal:   true,
	DistPATH:     "dist",
}

type SyncSetting struct {
	DBPath         string  `json:"db_path"`
	DefaultJobPath string  `json:"default_job_path"`
	StorePath      string  `json:"store_path"`
	Server         *Server `json:"server"`
}

var SyncSet = &SyncSetting{}

func LoadUserConfig(filePath string) error {
	var (
		content []byte
		err     error
	)

	if len(filePath) == 0 {
		return fmt.Errorf("配置文件名不能为空")
	}
	log.Infof("当前使用的配置文件为:%s", filePath)

	content, _ = file.ReadFromFile(filePath)
	err = json.Unmarshal(content, &SyncSet)
	if err != nil {
		return fmt.Errorf("导入用户配置出现错误: %w", err)
	}
	if SyncSet.Server == nil {
		SyncSet.Server = defaultServerSetting
	}
	log.Info("成功导入用户配置")
	return nil
}

func GetBindAddr(bind bool, port int) string {
	var prefix string
	if bind == false {
		prefix = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", prefix, port)
}

func GetDistPATH() string {
	if SyncSet.Server.DistPATH == "" {
		return "./dist/"
	} else {
		return SyncSet.Server.DistPATH
	}
}

func GetDefaultJobPath() string {
	if SyncSet.DefaultJobPath == "" {
		return "conf/job.json"
	} else {
		return SyncSet.DefaultJobPath
	}
}

func GetDBPath() string {
	if SyncSet.DBPath == "" {
		return "sync.db"
	} else {
		return SyncSet.DBPath
	}
}

func GetStorePath() string {
	if SyncSet.StorePath == "" {
		return "/data1"
	} else {
		return SyncSet.StorePath
	}
}
