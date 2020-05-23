package rsync

import (
	"testing"
	"time"
)

func TestRsyncCheck(t *testing.T) {
	hasRsync, info := CheckRsync()
	if !hasRsync {
		t.Errorf("rsync is not installed on this machine, stop.")
	}
	t.Logf("check rsync ok. version info: \n%s", info)
}

func TestExecCommandTimeout(t *testing.T) {
	goodTimeout := func(conf *Config) {
		t.Log("execute timeout ok.")
	}
	badTimeout := func(conf *Config) {
		t.Error("should not execute this timeout.")
	}
	conf := &Config{
		Command: "sleep",
		Args:    []string{"2"},
		Timeout: 1 * time.Second,
	}
	ExecCommand(conf, goodTimeout)
	conf.Timeout = 3
	ExecCommand(conf, badTimeout)
}
