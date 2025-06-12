package mutex

import (
	"context"
	"sync"
)

// Mutex is a drop-in replacement for the standard library's sync.Mutex. It
// offers the ability to cancel acquiring a lock by making use of Go's select
// statement. The zero value is safe to use. Satisfies sync.Locker.
type Mutex struct {
	// initMu is used to ensure that the mutex is initialized only once.
	initMu sync.Mutex

	// state must be a buffered channel with a length of 1. It is considered
	// locked when it has a length of 1. It is considered unlocked when it has a
	// length of 0.
	state chan struct{}
}

// Unlocks m. Panics if m is not locked. Calling Unlock on an unlocked mutex
// usually indicates a race condition.
func (m *Mutex) Unlock() {
	m.init()
	select {
	case <-m.state:
	default:
		panic("unlock of unlocked mutex")
	}
}

// SendLock locks m when sending to its returned channel. Do not close it.
// Calling SendLock while holding the lock will deadlock.
func (m *Mutex) SendLock() chan<- struct{} {
	m.init()
	return m.state
}

// Lock locks m. If the lock is already in use, the calling goroutine blocks
// until the mutex is available.
func (m *Mutex) Lock() {
	m.SendLock() <- struct{}{}
}

// TryLock tries to lock m and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare, and use of
// TryLock is often a sign of a deeper problem in a particular use of mutexes.
func (m *Mutex) TryLock() bool {
	select {
	case m.SendLock() <- struct{}{}:
		return true
	default:
		return false
	}
}

// GetLock locks m or returns ctx's error.
func (m *Mutex) GetLock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.SendLock() <- struct{}{}:
		return nil
	}
}

// init is used to fix a zero value Mutex.
func (m *Mutex) init() {
	m.initMu.Lock()
	defer m.initMu.Unlock()
	if m.state == nil {
		m.state = make(chan struct{}, 1)
	}
}

// IsLocked is mostly used for testing purposes. While the method is safe for
// concurrent use, the returned value should no be used as a synchonization
// mechanism as m's state may change before the returned value is evaluated.
func (m *Mutex) IsLocked() bool {
	return len(m.state) == 1
}
