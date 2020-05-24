package manager

// InitTaskKeys init task keys in cahce
func InitTaskKeys(keys []string) {

}

// StartTask notify manager a task had started executing
func StartTask(key string) error {
	return nil
}

// ExitTask notify manager a task had been executed
// it also record task information in history
func ExitTask(key string, err error) error {
	return nil
}
