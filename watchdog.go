package timeutil

import (
	"sync"
	"time"
)

// Watchdog executes a callback after the configured period without a reset.
// The timer starts immediately after creation. Reset and Stop are safe to call concurrently.
type Watchdog struct {
	mu       sync.Mutex
	timeout  time.Duration
	deadline time.Time
	timer    *time.Timer
	callback func()
	stopped  bool
}

// NewWatchdog creates a watchdog timer and starts it immediately.
func NewWatchdog(timeout time.Duration, callback func()) *Watchdog {
	t := &Watchdog{
		timeout:  timeout,
		callback: callback,
	}
	t.Reset()
	return t
}

// Reset restarts the timer from the current time.
func (t *Watchdog) Reset() {
	if t == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped || t.timeout <= 0 {
		return
	}

	t.deadline = time.Now().Add(t.timeout)
	if t.timer == nil {
		t.timer = time.AfterFunc(t.timeout, t.expire)
		return
	}

	t.timer.Reset(t.timeout)
}

// Stop permanently stops the timer. Subsequent Reset calls are ignored.
func (t *Watchdog) Stop() {
	if t == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		return
	}

	t.stopped = true
	t.deadline = time.Time{}
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
}

func (t *Watchdog) expire() {
	t.mu.Lock()
	if t.stopped || t.timer == nil {
		t.mu.Unlock()
		return
	}

	if remaining := time.Until(t.deadline); remaining > 0 {
		t.timer.Reset(remaining)
		t.mu.Unlock()
		return
	}

	t.stopped = true
	t.timer = nil
	callback := t.callback
	t.mu.Unlock()

	if callback != nil {
		callback()
	}
}
