package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type TaskStatus int

// History is the sync history of a task
type History struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Info      string    `json:"info"`
}

// task status enum
const (
	UNKNOWN TaskStatus = -1
	SLEEP   TaskStatus = 0
	STARTED TaskStatus = 1
	SUCC    TaskStatus = 2
	FAILED  TaskStatus = 3
)

const (
	statusPrefix      = "status:"
	timePrefix        = "time:"
	syncHistoryPrefix = "sync_history:"
	historyBucket     = "sync_history_bucket"
)

var cache = gocache.New(gocache.NoExpiration, 0)

// InitTaskKeys init task keys in cahce
func InitTaskKeys(keys []string) {
	for _, v := range keys {
		cache.Set(statusPrefix+v, SLEEP, gocache.NoExpiration)
	}
}

// StartTask notify manager a task had started executing
func StartTask(key string) error {
	status, ok := cache.Get(statusPrefix + key)
	if !ok {
		return fmt.Errorf("invalid task key: %s", key)
	}
	if status == STARTED {
		return fmt.Errorf("task %s had started, can not start again", key)
	}
	cache.Set(statusPrefix+key, STARTED, gocache.NoExpiration)
	cache.Set(timePrefix+key, time.Now(), gocache.NoExpiration)
	return nil
}

// ExitTask notify manager a task had been executed
// it also record task information in history
func ExitTask(key string, err error) error {
	status, ok := cache.Get(statusPrefix + key)
	if !ok {
		return fmt.Errorf("invalid task key: %s", key)
	}
	if status != STARTED {
		return fmt.Errorf("task %s status is %d, can not exit", key, status)
	}

	syncInfo := "ok"
	if err != nil {
		cache.Set(statusPrefix+key, FAILED, gocache.NoExpiration)
		syncInfo = err.Error()
	} else {
		cache.Set(statusPrefix+key, SUCC, gocache.NoExpiration)
	}
	startTime, _ := cache.Get(timePrefix + key)
	reocrdHistory(&History{
		Name:      key,
		StartTime: startTime.(time.Time),
		EndTime:   time.Now(),
		Info:      syncInfo,
	})
	return nil
}

// GetTaskStatus get task status from cache
func GetTaskStatus(key string) TaskStatus {
	ret, ok := cache.Get(statusPrefix + key)
	if !ok {
		return UNKNOWN
	}
	return ret.(TaskStatus)
}

// GetHistory returns all sync history about specified task
func GetHistory(taskName string) []*History {
	db, err := bolt.Open("sync.db", 0600, nil)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer db.Close()

	ret := make([]*History, 0)
	db.View(func(tx *bolt.Tx) error {
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
	return ret
}

func reocrdHistory(record *History) {
	db, err := bolt.Open("sync.db", 0600, nil)
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
