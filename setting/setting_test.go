package setting

import (
	"encoding/json"
	"github.com/PointStoneTeam/cugsync/manager"
	"github.com/PointStoneTeam/cugsync/rsync"
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	filePath := "../conf/config.json"
	LoadUserConfig(filePath)
	t.Log(SyncSet)
	t.Log(*SyncSet.Server)
}

func TestGetBindAddr(t *testing.T) {
	t.Log(GetBindAddr(false, 8000))
}

func TestGetDefaultJob(t *testing.T) {
	filePath := "../conf/job.json"
	if jobList, err := GetDefaultJob(filePath); err != nil {
		t.Log(err)
	} else {
		for i, j := range *jobList {
			t.Logf("job %d: %v", i, j)
		}
	}
}

func TestGetDefaultJob2(t *testing.T) {
	var list []manager.UnCreatedJob
	list = append(list, manager.UnCreatedJob{
		Name:   "job1",
		Spec:   "*/1 * * * *",
		Config: rsync.InitConfig("./test_rsync", "/eclipse/swtchart/releases/0.12.0/release/", "mirrors.tuna.tsinghua.edu.cn", []string{"-avz", "--delete"}),
	})
	str, _ := json.Marshal(list)
	t.Log(string(str))
}
