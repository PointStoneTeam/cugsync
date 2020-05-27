package setting

import (
	"encoding/json"
	"fmt"
	"github.com/PointStoneTeam/cugsync/manager"
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
	Server *Server `json:"server"`
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
	return SyncSet.Server.DistPATH
}

// resolve default job config
func GetDefaultJob(filePath string) (*[]manager.UnCreatedJob, error) {
	var (
		content      []byte
		err          error
		unCreatedJob *[]manager.UnCreatedJob
	)

	if len(filePath) == 0 {
		return nil, fmt.Errorf("默认任务配置文件名不能为空")
	}
	log.Infof("当前使用的任务计划配置文件为:%s", filePath)

	content, _ = file.ReadFromFile(filePath)
	err = json.Unmarshal(content, &unCreatedJob)
	if err != nil {
		return nil, fmt.Errorf("导入任务计划配置出现错误: %w", err)
	}
	log.Info("成功导入默认任务配置")
	return unCreatedJob, nil
}
