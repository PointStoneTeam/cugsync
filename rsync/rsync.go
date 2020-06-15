package rsync

import (
	"bytes"
	"github.com/PointStoneTeam/cugsync/pkg/file"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
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
	var err error
	// checkPath exist,if no folder then create
	if _, err = file.PathExists(conf.LocalDir); err != nil {
		err = os.Mkdir(conf.LocalDir, 755)
		if err != nil {
			log.WithFields(log.Fields{
				"err":      err,
				"localDir": conf.LocalDir,
			}).Info("Mkdir error")
			return err
		}
	}

	args := append(conf.Args, "rsync://"+conf.Upstream+conf.RemoteDir, "./")
	rsyncCmd := exec.Command(conf.Command, args...)
	rsyncCmd.Dir = conf.LocalDir

	// open a file and redirect std
	name := strings.Split(conf.LocalDir, "/")
	logName := name[len(name)-1] + ".log"
	os.Remove(logName)
	f, _ := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	//指定输出位置
	rsyncCmd.Stderr = f
	rsyncCmd.Stdout = f
	// run the command
	if err := rsyncCmd.Run(); err != nil {
		return err
	}
	return nil
}
