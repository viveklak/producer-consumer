## Producer-Consumer

This is a simple demonstration of a producer/consumer model managed by a self-managed job queue. 
When new jobs are added to the queue and there is a queueing delay of 1 second or more because there 
aren't enough active workers/consumers to handle the job, a new worker go routine is spawned. 
Similarly if the worker is found to be idle for a certain time (currently hard coded to be 5 seconds),
the worker shutsdown.

Simply run the main program:
```
  go run main.go
```
This will launch 5 producers. Typing `s` followed by `enter` will print a summary of current producers and consumers. Play around with adding new producers (`p` + `enter`) or killing existing producers (`k` + `enter`) to see what the state of the producers and consumers looks like after each of these operations.