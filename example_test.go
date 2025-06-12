package mutex_test

import (
	"context"
	"fmt"
	"time"

	"github.com/jakobii/mutex"
)

// Demonstates a few ways that Mutex can be used.
func ExampleMutex() {
	var mu mutex.Mutex

	// Lock
	mu.Lock()
	fmt.Println("A drop-in replacement for sync.Mutex.")
	mu.Unlock()

	// TryLock
	if mu.TryLock() {
		fmt.Println("Why not give it a try?")
		mu.Unlock()
	}

	// SendLock
	select {
	case mu.SendLock() <- struct{}{}:
		fmt.Println("This mutex plays nice with select statements.")
		mu.Unlock()
	default:
		return
	}

	// GetLock
	if err := mu.GetLock(context.Background()); err != nil {
		return
	}
	fmt.Println("Which means it can easily respect contexts.")
	mu.Unlock()

	// Output:
	// A drop-in replacement for sync.Mutex.
	// Why not give it a try?
	// This mutex plays nice with select statements.
	// Which means it can easily respect contexts.
}

func ExampleMutex_GetLock() {
	var mu mutex.Mutex

	// Imagine some other process has the lock.
	mu.Lock()

	// Attempt to acquire the lock with a context. This timeout is short for the
	// sake of the completion of the example. In practice this might be a larger
	// duration for detecting deadlocks, or a context tied to some other process
	// cancelation.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	if err := mu.GetLock(ctx); err != nil {
		fmt.Println("Failed to acquire lock:", err)
		return
	}
	defer mu.Unlock()

	// Output:
	// Failed to acquire lock: context deadline exceeded
}

func ExampleMutex_SendLock() {
	var mu mutex.Mutex

	// The process that sends on the mutex channel is the one that obtains the
	// lock.
	select {
	case <-time.After(time.Millisecond):
		fmt.Println("Failed to acquire lock: timeout")
		return
	case mu.SendLock() <- struct{}{}:
		defer mu.Unlock()
		// ..do things while holding the lock...
		fmt.Println("Successfully acquired lock.")
	}

	// Output:
	// Successfully acquired lock.
}
