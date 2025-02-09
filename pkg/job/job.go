package job

import (
	"context"
	"log"
	"time"
)

type Job struct {
	JobName    string
	BusyTime   time.Duration
	SubmitTime time.Time
	StartTime  time.Time
	EndTime    time.Time
}

func (j *Job) Run(ctx context.Context, workerName string) {
	log.Printf("Starting job: %+v on worker: %q", j, workerName)
	j.StartTime = time.Now()
	tick := time.Tick(j.BusyTime)
	select {
	case <-ctx.Done():
	case <-tick:
	}
	j.EndTime = time.Now()
	log.Printf("Completed job: %+v on worker: %q\n", j, workerName)
}
