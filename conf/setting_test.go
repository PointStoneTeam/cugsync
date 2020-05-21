package conf

import (
	"fmt"
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	filePath := "../config.json"
	LoadUserConfig(filePath)
	fmt.Println(SyncSet)
	fmt.Println(*SyncSet.Server)
}

func TestGetBindAddr(t *testing.T) {
	fmt.Println(GetBindAddr(false, 8000))
}
