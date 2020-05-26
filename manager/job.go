package manager

import (
	"fmt"
	"github.com/PointStoneTeam/cugsync/cron"
	"github.com/PointStoneTeam/cugsync/rsync"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

// task start status enum
type TaskStart int

const (
	Create TaskStart = 1
	Start  TaskStart = 2
	Stop   TaskStart = 3
)

// Define of Job
type Job struct {
	Name             string         `json:"name"`
	Spec             string         `json:"spec"`   // cron expression
	Config           *rsync.Config  `json:"config"` // rsync config
	StartTime        time.Time      `json:"start_time"`
	EndTime          time.Time      `json:"end_time"`
	Status           TaskStart      `json:"status"` // Create | Start | Stop
	LatestSyncStatus SyncTaskStatus `json:"latest_sync_status"`
	Shut             chan int       `json:"shut"` // use chan to stop job
}

// implement Run() interface to start rsync job
func (this Job) Run() {

	// start rsync job, record rsync history
	if err := rsync.ExecCommand(this.Config); err != nil {
		recordHistory(&History{
			Name:      this.Name,
			StartTime: this.StartTime,
			EndTime:   time.Now(),
			Info:      err.Error(),
		})
	} else {
		recordHistory(&History{
			Name:      this.Name,
			StartTime: this.StartTime,
			EndTime:   time.Now(),
			Info:      "rsync complete",
		})
	}

}

// CreateJob : create job, add to cache
// if job exist, this operation will replace job
func CreateJob(name, spec string, config *rsync.Config) {
	// init job status
	job := &Job{
		Name:             name,
		Spec:             spec,
		Config:           config,
		StartTime:        time.Now(),
		Status:           Create,
		Shut:             make(chan int),
		LatestSyncStatus: UNKNOWN,
	}

	cache.Set(jobPrefix+job.Name, job, gocache.NoExpiration)
}

// GetJob : get job from cache by name
func GetJob(name string) (*Job, error) {
	ret, ok := cache.Get(jobPrefix + name)
	if !ok {
		log.WithFields(log.Fields{
			"name": name,
		}).Info("未找到任务计划")
		return nil, fmt.Errorf("未找到任务计划 %s", name)
	}
	return ret.(*Job), nil
}

// StartJob :
func StartJob(name string) error {
	var (
		j   *Job
		err error
	)
	if j, err = GetJob(name); err != nil {
		return err
	}

	// start job
	if j.Status == Start {
		return fmt.Errorf("job: %s is started.", j.Name)
	}
	go cron.StartJob(j.Spec, j, j.Shut)
	j.Status = Start
	return nil
}

// StopJob
func StopJob(name string) error {
	var (
		j   *Job
		err error
	)
	if j, err = GetJob(name); err != nil {
		return err
	}

	// shut job
	if j.Shut == nil {
		return fmt.Errorf("job: %s shut error", j.Name)
	} else if j.Status == Stop {
		return fmt.Errorf("job: %s is started.", j.Name)
	}
	cron.StopJob(j.Shut)
	j.Status = Stop
	return nil
}
