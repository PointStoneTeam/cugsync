package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PointStoneTeam/cugsync/setting"
	"time"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

// History is the sync history of a task
type History struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Info      string    `json:"info"`
}

// sync task status enum
type SyncTaskStatus int

const (
	UNKNOWN SyncTaskStatus = -1
	SLEEP   SyncTaskStatus = 0
	STARTED SyncTaskStatus = 1
	SUCC    SyncTaskStatus = 2
	FAILED  SyncTaskStatus = 3
)

const (
	jobPrefix         = "job:"
	timePrefix        = "time:"
	syncHistoryPrefix = "sync_history:"
	historyBucket     = "sync_history_bucket"
)

// GetHistory returns all sync history about specified task
func GetHistory(taskName string) ([]*History, error) {
	db, err := bolt.Open(setting.GetDBPath(), 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ret := make([]*History, 0)
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(historyBucket)).Cursor()
		if c == nil {
			return fmt.Errorf("bucket is not init")
		}
		keyPrefix := fmt.Sprintf("%s%s", syncHistoryPrefix, taskName)
		for k, v := c.Seek([]byte(keyPrefix)); k != nil && bytes.HasPrefix(k, []byte(keyPrefix)); k, v = c.Next() {
			item := new(History)
			if err := json.Unmarshal(v, item); err != nil {
				return err
			}
			ret = append(ret, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func RecordHistory(record *History) {
	db, err := bolt.Open(setting.GetDBPath(), 0600, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer db.Close()

	// convert record to text
	infoBytes, err := json.Marshal(record)
	if err != nil {
		log.Error(err)
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucket))
		if b == nil {
			tx.CreateBucket([]byte(historyBucket))
			b = tx.Bucket([]byte(historyBucket))
		}
		key := fmt.Sprintf("%s%s%v", syncHistoryPrefix, record.Name, time.Now().Unix())
		b.Put([]byte(key), infoBytes)
		return nil
	})
	if err != nil {
		log.Error(err)
	}
}
