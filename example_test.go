package mutex_test

import (
	"context"
	"fmt"

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

	// WaitLock
	select {
	case mu.WaitLock() <- struct{}{}:
		fmt.Println("This mutex plays nice with select statements.")
		mu.Unlock()
	default:
		return
	}

	// LockCtx
	if err := mu.LockCtx(context.Background()); err != nil {
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
