package worker

import (
	"container/list"
	"context"
	"fmt"
	"github.com/viveklak/producer-consumer/pkg/active"
	"github.com/viveklak/producer-consumer/pkg/job"
	"time"
)

type Worker struct {
	ctx             context.Context
	cancelFunc      context.CancelFunc
	name            string
	startTime       time.Time
	terminationTime time.Time
	lastActiveTime  time.Time
	queue           <-chan *job.Job
	completedJobs   uint64
	activeWorkers   *active.ActiveList[*Worker]
	element         *list.Element
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) String() string {
	return fmt.Sprintf("Worker %q: Completed jobs: %d - lastActiveTime: %v - elapsed: %v",
		w.name, w.completedJobs, w.lastActiveTime, time.Since(w.lastActiveTime))
}

func (w *Worker) Shutdown() {
	if w.terminationTime != (time.Time{}) {
		return
	}
	fmt.Printf("[%s] Triggering shutdown\n", w.name)
	w.terminationTime = time.Now()
	w.activeWorkers.Remove(w.element)
	w.cancelFunc()
}

func (w *Worker) Run() error {
	for {
		tick := time.Tick(5 * time.Second)
		select {
		case j := <-w.queue:
			w.lastActiveTime = time.Now()
			j.Run(w.ctx, w.name)
			w.completedJobs++
		case <-tick:
			w.Shutdown()
			return nil
		case <-w.ctx.Done():
			w.Shutdown()
			return nil
		}
	}
}

func NewWorker(ctx context.Context, workerID uint64, activeWorkers *active.ActiveList[*Worker], queue <-chan *job.Job) *Worker {
	newCtx, cancelFunc := context.WithCancel(ctx)
	w := &Worker{
		ctx:           newCtx,
		cancelFunc:    cancelFunc,
		name:          fmt.Sprintf("worker_%d", workerID),
		startTime:     time.Now(),
		queue:         queue,
		activeWorkers: activeWorkers,
	}
	e := activeWorkers.Insert(w)
	w.element = e
	go w.Run()
	return w
}
