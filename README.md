# mutex

Go mutex that plays nice with the select statement.

```
go get github.com/jakobii/mutex
```

The `mutex.Mutex` works with `select` statements and as a result can be used with `context.Context`s.

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

Work like `sync.Mutex` as well.

```go
var mu mutex.Mutex
mu.Lock()
defer mu.Unlock()
```