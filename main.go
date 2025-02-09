package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/viveklak/producer-consumer/pkg/active"
	"github.com/viveklak/producer-consumer/pkg/producer"
	"github.com/viveklak/producer-consumer/pkg/queue"
	"github.com/viveklak/producer-consumer/pkg/worker"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	fmt.Println("Hello world!")

	ctx, cancelFunc := context.WithCancel(context.Background())
	rch := make(chan rune, 1)
	go readKey(ctx, rch)

	wg := &sync.WaitGroup{}
	var signalChan chan (os.Signal) = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	activeWorkers := active.NewActiveList[*worker.Worker]()
	jq := queue.NewJobQueue(activeWorkers)

	activeProducers := active.NewActiveList[*producer.Producer]()
	var i int
	for i = 0; i < 5; i++ {
		p := producer.NewProducer(ctx, fmt.Sprintf("producer_%d", i), activeProducers, wg)
		go p.Run(jq)
	}

	bail := func(waitGroup *sync.WaitGroup) {
		cancelFunc()
		waitGroup.Wait()
	}

	for {
		select {
		case <-signalChan:
			bail(wg)
			return
		case k := <-rch:
			if k == 'q' {
				fmt.Println("**** Shutting down...")
				bail(wg)
				return
			}

			if k == 'p' {
				fmt.Println("**** Adding new producer...")
				p := producer.NewProducer(ctx, fmt.Sprintf("producer_%d", i), activeProducers, wg)
				i++
				go p.Run(jq)
			}

			if k == 'k' {
				fmt.Println("**** Removing producer...")
				p, err := activeProducers.DropFront()
				if err != nil {
					panic(err)
				} else {
					log.Printf("Shutting down [%v]\n", p)
					p.Shutdown()
				}
			}

			if k == 's' {
				fmt.Println("**** Printing summary of producers...")
				activeProducers.Print()
				fmt.Println("**** Printing summary of consumers/workers...")
				activeWorkers.Print()
			}
		}
	}
}

func readKey(ctx context.Context, input chan rune) {
	var reader = bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		char, _, err := reader.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		input <- char
	}
}
