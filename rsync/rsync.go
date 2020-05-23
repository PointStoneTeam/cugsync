package rsync

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const RSYNC_CMD = "rsync"

// Config represents config for each sync task
type Config struct {
	Command   string // set Command as `rsync` on default
	LocalDir  string
	RemoteDir string
	Args      []string
	Timeout   time.Duration
}

// TimeoutCallback represents timeout call back function for rsync command
type TimeoutCallback func(*Config)

// InitConfig init and return a rsync task config
func InitConfig(localDir, remoteDir string, args []string, timeout time.Duration) *Config {
	return &Config{
		Command:   RSYNC_CMD,
		LocalDir:  localDir,
		RemoteDir: remoteDir,
		Args:      args,
		Timeout:   timeout,
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
// it will call `callback` and exit function if timeout
func ExecCommand(conf *Config, callback TimeoutCallback) {
	ch := make(chan bool)

	go func() {
		rsyncCmd := exec.Command(conf.Command, conf.Args...)
		rsyncCmd.Run()
		ch <- true
	}()

	select {
	case <-ch:

	case <-time.After(conf.Timeout):
		go callback(conf)
	}
}

// DefaultTimeoutCallback print info about task in default
func DefaultTimeoutCallback(conf *Config) {
	log.Errorf(`Sync task cost %s timeout. Full command %s %s.
	\nLocal dir: %s\nRemote Dir: %s\n`, conf.Timeout.String(),
		conf.Command, strings.Join(conf.Args, " "), conf.LocalDir, conf.RemoteDir)
}
