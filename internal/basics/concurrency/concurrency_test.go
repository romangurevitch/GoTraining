package concurrency

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

type counter struct {
	counter int
}

// Bad example how to use go functions - run to see the results.
// Not thread safe, run with -race to find the issue.
func TestSimplestGoFunc(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	go func() {
		for i := 0; i < 1000; i++ {
			c.counter++
		}
	}()

	// What will be printed?
	fmt.Println(c.counter)
	fmt.Println("Terminating...")
}

// Unexpected results since it is not thread safe, run with -race to find the issue.
func TestGoFuncWithWaitGroup(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c.counter++
		}()
	}

	wg.Wait()

	// What will be printed?
	fmt.Println(c.counter)
	fmt.Println("Terminating...")
}

// Good example with expected results and thread safe.
func TestGoFuncWithWaitGroupMutex(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			lock.Lock()
			c.counter++
			lock.Unlock()
		}()
	}

	wg.Wait()

	// What will be printed?
	fmt.Println(c.counter)
	fmt.Println("Terminating...")
}

// Bad endless go func example, terminating the run will terminate without printing the results.
func TestGoFuncEndless(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	wg.Add(1)
	go func() {
		for {
			lock.Lock()
			c.counter++
			lock.Unlock()
		}
	}()

	wg.Wait()

	// What will be printed?
	fmt.Println(c.counter)
	fmt.Println("Terminating...")
}

// Using OS signals to catch termination signal to print out counter results.
// Good example how to make sure resources are closed when terminating running processes.
func TestGoFuncEndlessWithChannel(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	lock := sync.Mutex{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			lock.Lock()
			c.counter++
			lock.Unlock()
		}
	}()

	fmt.Println("Waiting...")
	<-sigs

	// What will be printed?
	lock.Lock()
	fmt.Println(c.counter)
	lock.Unlock()

	fmt.Println("Terminating...")
}

// Producer example, to show how signals behave when channel buffer is full without consumer.
// Run to see the results.
func TestChannelProducer(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	lock := sync.Mutex{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	numChan := make(chan int, 10)

	// Producer
	go func() {
		for {
			lock.Lock()
			c.counter++

			numChan <- c.counter
			fmt.Println(c.counter)
			lock.Unlock()
		}
	}()

	fmt.Println("Waiting...")
	s := <-sigs

	fmt.Println(s)
	fmt.Println("Terminating...")
}

// Simple producer consumer.
func TestChannelProducerConsumer(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	lock := sync.Mutex{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	numChan := make(chan int)

	// Producer
	go func() {
		for {
			lock.Lock()
			c.counter++

			numChan <- c.counter
			fmt.Printf("producer %d\n", c.counter)
			lock.Unlock()
		}
	}()

	// Consumer
	go func() {
		for {
			counter := <-numChan
			fmt.Printf("\t\t\t\tconsumer_0 %d\n", counter)
			time.Sleep(time.Second)
		}
	}()

	fmt.Println("Processing...")
	s := <-sigs

	fmt.Println(s)
	fmt.Println("Terminating...")
}

// Single producer multiple consumers.
func TestChannelProducerMultipleConsumers(t *testing.T) {
	t.Skip("Skipping: example test, comment out to run manually...")

	fmt.Println("Starting the test...")
	c := &counter{}
	lock := sync.Mutex{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	numChan := make(chan int)

	// Producer
	go func() {
		for {
			lock.Lock()
			c.counter++

			numChan <- c.counter
			fmt.Printf("producer %d\n", c.counter)
			lock.Unlock()
		}
	}()

	// Consumers
	for i := 0; i < 3; i++ {
		consumerNum := i
		go func() {
			for {
				counter := <-numChan
				fmt.Printf("\t\t\t\tconsumer_%d %d\n", consumerNum, counter)
				time.Sleep(time.Second * 3)
			}
		}()
	}

	fmt.Println("Processing...")
	s := <-sigs

	fmt.Println(s)
	fmt.Println("Terminating...")
}

func TestCloseChannel(t *testing.T) {
	simpleChannel := make(chan int)
	close(simpleChannel)

	_, ok := <-simpleChannel
	if ok {
		t.Fatal("something is wrong")
	}
}

func TestWriteToClosedChannel(t *testing.T) {
	simpleChannel := make(chan int)
	close(simpleChannel)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered", r)
		}
	}()

	simpleChannel <- 1
}

func TestCloseChannelForLoop(t *testing.T) {
	simpleChannel := make(chan int)
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond) // nolint

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-simpleChannel:
				fmt.Println("value received:", val)
			}
		}
	}()

	close(simpleChannel)
	<-ctx.Done()
}

func TestCloseMultipleChannel(t *testing.T) {
	simpleChannel := make(chan int)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered", r)
		}
	}()

	close(simpleChannel)
	close(simpleChannel)
}

func TestCtxCloseChannel(t *testing.T) {
	simpleChannel := make(chan int)
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered", r)
			}
		}()

		for { // nolint
			select {
			case <-ctx.Done():
				fmt.Println("attempting to close")
				close(simpleChannel)
			}
		}
	}()

	cancelFunc()
	<-simpleChannel
}
