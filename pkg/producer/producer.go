package producer

import (
	"container/list"
	"context"
	"fmt"
	"github.com/viveklak/producer-consumer/pkg/active"
	"github.com/viveklak/producer-consumer/pkg/queue"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Producer struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	name         string
	rand         *rand.Rand
	sleepTime    time.Duration
	queueingTime time.Duration
	wg           *sync.WaitGroup
	element      *list.Element
}

// Implement stringer
func (p *Producer) String() string {
	// Include information about sleepTime, queueingTime and name
	return fmt.Sprintf("Producer %q: sleepTime: %v, queueingTime: %v", p.name, p.sleepTime, p.queueingTime)
}

func (p *Producer) Name() string {
	return p.name
}

func NewProducer(ctx context.Context, name string, activeProducers *active.ActiveList[*Producer], waitGroup *sync.WaitGroup) *Producer {
	newCtx, cancelFunc := context.WithCancel(ctx)

	p := &Producer{
		ctx:        newCtx,
		cancelFunc: cancelFunc,
		name:       name,
		rand:       rand.New(rand.NewSource(time.Now().Unix())),
		sleepTime:  time.Duration(rand.Int63n(5000)+1) * time.Millisecond,
		wg:         waitGroup,
	}
	el := activeProducers.Insert(p)
	p.element = el
	return p
}

func (p *Producer) Shutdown() {
	p.cancelFunc()
}

func (p *Producer) Run(q *queue.JobQueue) error {
	var ctr int64

	p.wg.Add(1)
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return nil
		case <-time.Tick(p.sleepTime):
			jobTime := p.rand.Int63n(100) + 1
			jobName := fmt.Sprintf("%s_%d", p.name, ctr)
			start := time.Now()
			log.Printf("Submitting new job: %q", jobName)
			q.SubmitJob(p.ctx, jobName, time.Duration(jobTime)*time.Second)
			p.queueingTime += time.Since(start)
			ctr += 1
		}
	}
}
