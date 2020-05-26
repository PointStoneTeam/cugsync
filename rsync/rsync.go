package rsync

import (
	"bytes"
	"github.com/PointStoneTeam/cugsync/pkg/file"
	"os"
	"os/exec"
)

const RSYNC_CMD = "rsync"

// Config represents config for each sync task
type Config struct {
	Command   string // set Command as `rsync` on default
	LocalDir  string
	RemoteDir string
	Upstream  string
	Args      []string
}

// InitConfig init and return a rsync task config
func InitConfig(localDir, remoteDir, upstream string, args []string) *Config {
	return &Config{
		Command:   RSYNC_CMD,
		LocalDir:  localDir,
		RemoteDir: remoteDir,
		Upstream:  upstream,
		Args:      args,
	}
}

// CheckRsync checks if rsync had been installed
// and returns the version info about rsync
func CheckRsync() (hasRsync bool, info string) {
	rsyncCmd := exec.Command(RSYNC_CMD, "--version")
	var out bytes.Buffer
	rsyncCmd.Stdout = &out
	if err := rsyncCmd.Run(); err != nil {
		return false, err.Error()
	}
	return true, out.String()
}

// ExecCommand start calling rsync command with specified config
func ExecCommand(conf *Config) error {
	// checkPath exist,if no folder then create
	if _, err := file.PathExists(conf.LocalDir); err != nil {
		os.MkdirAll(conf.LocalDir, 666)
	}

	args := append(conf.Args, "rsync://"+conf.Upstream+conf.RemoteDir, "./")
	rsyncCmd := exec.Command(conf.Command, args...)
	rsyncCmd.Dir = conf.LocalDir
	if err := rsyncCmd.Run(); err != nil {
		return err
	}
	return nil
}
