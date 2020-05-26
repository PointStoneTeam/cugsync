package manager

import (
	"github.com/PointStoneTeam/cugsync/rsync"
	"os"
	"testing"
	"time"
)

func TestCreateJob(t *testing.T) {
	CreateJob("testJob", "* * * * *", rsync.InitConfig("./test_rsync",
		"/eclipse/swtchart/releases/0.12.0/release/",
		"mirrors.tuna.tsinghua.edu.cn",
		[]string{"-avz", "--delete"}))
	defer os.Remove("sync.db")

	var (
		j   *Job
		err error
	)

	if j, err = GetJob("testJob"); err != nil {
		t.Error(err)
	}
	t.Logf("job: %v", j)
}

func TestStartJob(t *testing.T) {
	jobName := "testJob"
	// every minute do this job
	CreateJob(jobName, "*/1 * * * *", rsync.InitConfig("./test_rsync/",
		"/eclipse/swtchart/releases/0.12.0/release/",
		"mirrors.tuna.tsinghua.edu.cn",
		[]string{"-avz", "--delete"}))
	defer os.Remove("sync.db")

	var (
		j   *Job
		err error
	)

	if j, err = GetJob("testJob"); err != nil {
		t.Error(err)
	}
	t.Logf("job: %v", j)

	StartJob(jobName)
	time.Sleep(time.Minute)
	historyList, err := GetHistory("testJob")
	if err != nil {
		t.Errorf("GetHistory error: %s", err.Error())
	}
	for index, history := range historyList {
		t.Logf("index %d : %v", index, history)
	}
	StopJob(jobName)
	t.Logf("job.Name : %s, job.Status: %d", j.Name, j.Status)
}
