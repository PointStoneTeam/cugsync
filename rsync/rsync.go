package rsync

import (
	"bytes"
	"os/exec"
)

const RSYNC_CMD = "rsync"

// Config represents config for each sync task
type Config struct {
	Command   string // set Command as `rsync` on default
	LocalDir  string
	RemoteDir string
	Args      []string
}

// TimeoutCallback represents timeout call back function for rsync command
type TimeoutCallback func(*Config)

// ErrorCallback represents callback function when
type ErrorCallback func(*Config, error)

// InitConfig init and return a rsync task config
func InitConfig(localDir, remoteDir string, args []string) *Config {
	return &Config{
		Command:   RSYNC_CMD,
		LocalDir:  localDir,
		RemoteDir: remoteDir,
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
	rsyncCmd := exec.Command(conf.Command, conf.Args...)
	if err := rsyncCmd.Run(); err != nil {
		return err
	}
	return nil
}
