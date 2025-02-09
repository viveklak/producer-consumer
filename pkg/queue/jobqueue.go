package queue

import (
	"context"
	"github.com/viveklak/producer-consumer/pkg/active"
	"github.com/viveklak/producer-consumer/pkg/job"
	"github.com/viveklak/producer-consumer/pkg/worker"
	"time"
)

type JobQueue struct {
	queue         chan *job.Job
	workerID      uint64
	activeWorkers *active.ActiveList[*worker.Worker]
}

func NewJobQueue(activeWorkers *active.ActiveList[*worker.Worker]) *JobQueue {
	return &JobQueue{
		queue:         make(chan *job.Job),
		activeWorkers: activeWorkers,
	}
}

func (j *JobQueue) SubmitJob(ctx context.Context, jobName string, busyTime time.Duration) {
	for {
		tick := time.Tick(1 * time.Second)
		select {
		case j.queue <- &job.Job{
			JobName:    jobName,
			BusyTime:   busyTime,
			SubmitTime: time.Now(),
		}:
			return
		case <-tick:
			worker.NewWorker(ctx, j.workerID, j.activeWorkers, j.queue)
			j.workerID += 1
		case <-ctx.Done():
			return
		}
	}
}

func (j *JobQueue) GetQueue() <-chan *job.Job {
	return j.queue
}
