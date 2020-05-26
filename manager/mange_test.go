package manager

import "testing"

func TestInitTask(t *testing.T) {
	keys := []string{"ubuntu", "centos", "archlinux", "debian"}
	InitTaskKeys(keys)
	for _, v := range keys {
		if s := GetTaskStatus(v); s != SLEEP {
			t.Errorf("init tasks error, get status %v", s)
		}
	}
}

func TestModifyTaskStatus(t *testing.T) {
	keys := []string{"ubuntu", "centos", "archlinux", "debian"}
	InitTaskKeys(keys)
	if err := StartTask("ubuntu"); err != nil {
		t.Error(err)
	}
	if s := GetTaskStatus("ubuntu"); s != STARTED {
		t.Errorf("task ubuntu started but get status %v", s)
	}
	if err := ExitTask("ubuntu", nil); err != nil {
		t.Error(err)
	}
	// test if history reocrded
	history, err := GetHistory("ubuntu")
	if err != nil {
		t.Error(err)
	}
	if len(history) == 0 {
		t.Errorf("bad history record")
	}
	t.Logf("task %s started at %v ended at %v, info %s", history[0].Name, history[0].StartTime, history[0].EndTime, history[0].Info)
}
