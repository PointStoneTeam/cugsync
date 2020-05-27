package manager

import (
	"github.com/PointStoneTeam/cugsync/rsync"
	"os"
	"testing"
	"time"
)

func TestCreateJob(t *testing.T) {
	CreateJob(&UnCreatedJob{
		Name: "testJob",
		Spec: "* * * * *",
		Config: rsync.InitDefaultConfig("./test_rsync",
			"/eclipse/swtchart/releases/0.12.0/release/",
			"mirrors.tuna.tsinghua.edu.cn"),
	})
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
	CreateJob(&UnCreatedJob{
		Name: "testJob",
		Spec: "* * * * *",
		Config: rsync.InitDefaultConfig("./test_rsync",
			"/eclipse/swtchart/releases/0.12.0/release/",
			"mirrors.tuna.tsinghua.edu.cn"),
	})
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

func TestInitJobs(t *testing.T) {
	var list []UnCreatedJob
	list = append(list, UnCreatedJob{
		Name:   "job1",
		Spec:   "*/1 * * * *",
		Config: rsync.InitConfig("./test_rsync", "/eclipse/swtchart/releases/0.12.0/release/", "mirrors.tuna.tsinghua.edu.cn", []string{"-avz", "--delete"}),
	})
	InitJobs(list)
	if jobList, err := GetAllJobs(); err != nil {
		t.Errorf("getalljobs error: %s", err.Error())
	} else {
		for i, j := range jobList {
			t.Logf("job %d: %v", i, *j)
		}
	}
}
