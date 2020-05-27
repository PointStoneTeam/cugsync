package setting

import (
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
