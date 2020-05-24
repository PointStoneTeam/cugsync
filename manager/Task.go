package manager

import (
	"github.com/PointStoneTeam/cugsync/rsync"
	"time"
)

// 任务开始状态
const (
	start = iota
	stop
)

// Define of Task
type Task struct {
	Name          string        `json:"name"`
	Spec          string        `json:"spec"`   // cron expression
	Config        *rsync.Config `json:"config"` // rsync config
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	CurrentStatus int           `json:"current_status"` // start or stop
}
