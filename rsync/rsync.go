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
	Command   string   `json:"command"` // set Command as `rsync` on default
	LocalDir  string   `json:"local_dir"`
	RemoteDir string   `json:"remote_dir"`
	Upstream  string   `json:"upstream"`
	Args      []string `json:"args"`
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

// InitDefaultConfig init and return default args
func InitDefaultConfig(localDir, remoteDir, upstream string) *Config {
	return &Config{
		Command:   RSYNC_CMD,
		LocalDir:  localDir,
		RemoteDir: remoteDir,
		Upstream:  upstream,
		Args:      []string{"-avz", "--delete", "--ipv6", "--safe-links"},
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
// default args: -avz --delete --ipv6 --safe-links
func ExecCommand(conf *Config) error {
	// checkPath exist,if no folder then create
	if _, err := file.PathExists(conf.LocalDir); err != nil {
		os.Mkdir(conf.LocalDir, 755)
	}

	args := append(conf.Args, "rsync://"+conf.Upstream+conf.RemoteDir, "./")
	rsyncCmd := exec.Command(conf.Command, args...)
	rsyncCmd.Dir = conf.LocalDir
	if err := rsyncCmd.Run(); err != nil {
		return err
	}
	return nil
}
