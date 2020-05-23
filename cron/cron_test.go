package cron

import (
	"fmt"
	"testing"
	"time"
)

type TestJob struct {
}

func (this TestJob) Run() {
	i := 0
	for {
		fmt.Println(i)
		i++
		time.Sleep(time.Second)
	}
}

func TestStartJob(t *testing.T) {
	ch := make(chan int)
	spec := "*/1 * * * *"
	go StartJob(spec, TestJob{}, ch)
	time.Sleep(time.Minute * 3)
	StopJob(ch)
}
