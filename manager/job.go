package manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PointStoneTeam/cugsync/cron"
	"github.com/PointStoneTeam/cugsync/pkg/file"
	"github.com/PointStoneTeam/cugsync/rsync"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// task start status enum
type TaskStart int

const (
	Create TaskStart = 1
	Start  TaskStart = 2
	Stop   TaskStart = 3
)

var cache = gocache.New(gocache.NoExpiration, 0)

const JobList = "job_list"

func init() {
	var jobList []string
	cache.Set(JobList, jobList, gocache.NoExpiration)
}

// Define of Job
type Job struct {
	Name             string         `json:"name"`
	DisplayName      string         `json:"display_name"`
	Description      string         `json:"description"`
	Catalog          string         `json:"catalog"`
	Spec             string         `json:"spec"`   // cron expression
	Config           *rsync.Config  `json:"config"` // rsync config
	StartTime        time.Time      `json:"start_time"`
	EndTime          time.Time      `json:"end_time"`
	Status           TaskStart      `json:"status"` // Create | Start | Stop
	LatestSyncTime   time.Time      `json:"latest_sync_time"`
	LatestSyncStatus SyncTaskStatus `json:"latest_sync_status"`
	Shut             chan int       `json:"-"` // use chan to stop job
}

type UnCreatedJob struct {
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	Description string        `json:"description"`
	Catalog     string        `json:"catalog"`
	Spec        string        `json:"spec"`
	Config      *rsync.Config `json:"config"`
}

// implement Run() interface to start rsync job
func (j *Job) Run() {
	j.LatestSyncStatus = STARTED
	j.LatestSyncTime = time.Now()

	// start rsync job, record rsync history
	//log.Infof("start rsync job %s, upstream: %s, remoteDir: %s, localDir: %s, args: %v", j.Name, j.Config.Upstream, j.Config.RemoteDir, j.Config.LocalDir, j.Config.Args)
	log.Infof("start rsync job %s, spec: %s, config: %v", j.Name, j.Spec, j.Config)
	for retryCount := 3; retryCount >= 0; retryCount-- {
		// first time error, after a minute retry
		if err := rsync.ExecCommand(j.Config); err != nil && retryCount > 0 {
			log.WithFields(log.Fields{
				"retryCount": retryCount,
				"err":        err,
			}).Infof("error, rsync job %s, spec: %s, config: %v", j.Name, j.Spec, j.Config)
			<-time.After(5 * time.Minute)
		} else if err != nil {
			// job maybe failed
			j.LatestSyncStatus = FAILED
			RecordHistory(&History{
				Name:      j.Name,
				StartTime: j.StartTime,
				EndTime:   time.Now(),
				Info:      err.Error(),
			})
		} else {
			j.LatestSyncStatus = SUCC
			RecordHistory(&History{
				Name:      j.Name,
				StartTime: j.StartTime,
				EndTime:   time.Now(),
				Info:      "rsync complete",
			})
			break
		}
	}

}

// CreateJob : create job, add to cache
// if job exist, this operation will replace job
func CreateJob(j *UnCreatedJob) {
	// init job status
	job := &Job{
		Name:             j.Name,
		DisplayName:      j.DisplayName,
		Description:      j.Description,
		Catalog:          j.Catalog,
		Spec:             j.Spec,
		Config:           j.Config,
		StartTime:        time.Now(),
		Status:           Create,
		LatestSyncTime:   time.Now(),
		LatestSyncStatus: UNSTART,
		Shut:             make(chan int),
		// 使用指针才可以更新 job 状态
	}

	cache.Set(jobPrefix+job.Name, job, gocache.NoExpiration)
	ListAddJob(job.Name)
	log.Infof("create job: %s", j.Name)
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
	log.Infof("start job %s", j.Name)
	go cron.StartJob(j.Spec, j, j.Shut)
	j.Status = Start
	j.LatestSyncTime = time.Now()
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

// add jobName to cache list
func ListAddJob(jobName string) {
	if ret, found := cache.Get(JobList); found {
		jobList := ret.([]string)
		jobList = append(jobList, jobName)
		cache.Set(JobList, jobList, gocache.NoExpiration)
	}
}

func GetAllJobs() ([]*Job, error) {
	var (
		jobList []string
		reList  []*Job
	)
	if ret, found := cache.Get(JobList); found {
		jobList = ret.([]string)
	} else {
		return nil, fmt.Errorf("获取任务计划列表出错")
	}
	// get all from cache
	for _, jobName := range jobList {
		if ret, found := cache.Get(jobPrefix + jobName); found {
			reList = append(reList, ret.(*Job))
		}
	}
	return reList, nil
}

// InitJobs from conf,then create and start
func InitJobs(jList *[]UnCreatedJob) {
	for _, j := range *jList {
		CreateJob(&j)
		StartJob(j.Name)
	}
}

// resolve default job config
func GetDefaultJob(filePath string) (*[]UnCreatedJob, error) {
	var (
		content      []byte
		err          error
		unCreatedJob *[]UnCreatedJob
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
