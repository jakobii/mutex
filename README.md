# mutex

Go mutex that plays nice with the select statement.

```
go get github.com/jakobii/mutex
```

The `mutex.Mutex` works with Go's `select` statements and can be used with tools like `context.Context`s to cancel obtaining the lock. Sending `struct{}{}` is a common way of signaling. Here we use it to obtain the lock.

```go
var mu mutex.Mutex
select {
case <-ctx.Done():
	return ctx.Err()
case mu.WaitLock() <- struct{}{}:
	defer mu.Unlock()
	fmt.Println("obtained lock")
}
```

The `mutex.Mutex` feels like `sync.Mutex` as well.

```go
var mu mutex.Mutex
mu.Lock()
defer mu.Unlock()
```