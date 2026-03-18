package context

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/romangurevitch/go-training/internal/basics/entity"
)

// https://blog.golang.org/context
func TestContextForLogging(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, entity.New("key"), entity.New("external field val"))

	fmt.Printf("%s\n", ctx)
}

// Using WithCancel context
func TestContextWithCancel(t *testing.T) {
	fmt.Println("Starting the test...")
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 5)
		cancelFunc()
	}()

	fmt.Println("Waiting...")
	<-ctx.Done()

	// What will be printed?
	fmt.Println("Reason:", ctx.Err())
	fmt.Println("Terminating...")
}

// Using WithTimeout context
func TestContextWithCancelTimeout(t *testing.T) {
	fmt.Println("Starting the test...")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)

	go func() {
		time.Sleep(time.Second * 10)
		cancelFunc()
	}()

	fmt.Println("Waiting...")
	<-ctx.Done()

	// What will be printed?
	fmt.Println("Reason:", ctx.Err())
	fmt.Println("Terminating...")
}
