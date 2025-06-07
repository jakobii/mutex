# mutex

A Go mutex that is cancelable, channel-based so it plays nicely with the
`select` statement, and is a drop-in replacement for `sync.Mutex`.

```
go get github.com/jakobii/mutex
```

The `mutex.Mutex` can work just like `sync.Mutex`.

```go
var mu mutex.Mutex
mu.Lock()
defer mu.Unlock()
```

The standard library's `sync.Mutex` does not offer a way to cancel its `Lock()`
method while it is blocking to acquire the lock. This is usually fine when the
mutex is guarding a resource that does not take much time to access, like a
struct field. But when synchronizing longer-running processes, the need to
cancel work that is not needed anymore is frequently an issue.

If all you need is to cancel waiting on a lock with a context, this module's
`mutex.Mutex` has a convenient method for this.

```go
var mu mutex.Mutex
func MyWork(ctx context.Context) error {
	if err := mu.GetLock(ctx); err != nil {
		return fmt.Errorf("my work was canceled: %w", err)
	}
	defer mu.Unlock()
	fmt.Println("acquired lock")
}
```

The `mutex.Mutex` can work with Go's `select` statement. Sending `struct{}{}` is
a common way of signaling, and here we use it with `SendLock()` to acquire the
lock.

```go
var mu mutex.Mutex
func MyWork(ctx context.Context) error {
	select {
	case <-someCloseSignal:
		return errClosed
	case <-time.After(someMaxDuration):
		return errTimeout
	case <-ctx.Done():
		return fmt.Errorf("my work was canceled: %w", err)
	case mu.SendLock() <- struct{}{}:
		defer mu.Unlock()
		fmt.Println("acquired lock")
	}
}
```

See more [examples](./example_test.go).