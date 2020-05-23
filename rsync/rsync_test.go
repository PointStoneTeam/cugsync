package rsync

import (
	"testing"
)

func TestRsyncCheck(t *testing.T) {
	hasRsync, info := CheckRsync()
	if !hasRsync {
		t.Errorf("rsync is not installed on this machine, stop.")
	}
	t.Logf("check rsync ok. version info: \n%s", info)
}
