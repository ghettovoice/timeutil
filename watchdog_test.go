package timeutil_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ghettovoice/timeutil"
)

func ExampleWatchdog() {
	expired := make(chan struct{})

	// NewWatchdog starts a timer that fires after the configured inactivity
	// period unless it is reset or stopped.
	watchdog := timeutil.NewWatchdog(10*time.Millisecond, func() { close(expired) })
	defer watchdog.Stop()

	// Resetting the watchdog postpones its expiration.
	watchdog.Reset()
	time.Sleep(time.Millisecond)

	select {
	case <-expired:
		fmt.Println("expired")
	case <-time.After(5 * time.Millisecond):
		fmt.Println("still alive")
	}

	// Output:
	// still alive
}

func TestWatchdog_ExpiresAfterCreation(t *testing.T) {
	t.Parallel()

	expired := make(chan struct{})
	timer := timeutil.NewWatchdog(10*time.Millisecond, func() { close(expired) })
	defer timer.Stop()

	select {
	case <-expired:
	case <-time.After(time.Second):
		t.Fatal("timer did not expire")
	}
}

func TestWatchdog_ResetDelaysExpiration(t *testing.T) {
	t.Parallel()

	expired := make(chan struct{})
	timer := timeutil.NewWatchdog(30*time.Millisecond, func() { close(expired) })
	defer timer.Stop()

	time.Sleep(15 * time.Millisecond)
	timer.Reset()

	select {
	case <-expired:
		t.Fatal("timer expired before reset duration elapsed")
	case <-time.After(10 * time.Millisecond):
	}

	select {
	case <-expired:
	case <-time.After(time.Second):
		t.Fatal("timer did not expire after reset")
	}
}

func TestWatchdog_StopPreventsExpiration(t *testing.T) {
	t.Parallel()

	expired := make(chan struct{})
	timer := timeutil.NewWatchdog(10*time.Millisecond, func() { close(expired) })
	timer.Stop()
	timer.Reset()

	select {
	case <-expired:
		t.Fatal("stopped timer expired")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestWatchdog_ConcurrentReset(t *testing.T) {
	t.Parallel()

	var (
		mu      sync.Mutex
		expires int
		expired = make(chan struct{})
		timer   = timeutil.NewWatchdog(20*time.Millisecond, func() {
			mu.Lock()
			expires++
			mu.Unlock()
			close(expired)
		})
	)
	defer timer.Stop()

	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			for range 20 {
				timer.Reset()
			}
		})
	}
	wg.Wait()

	select {
	case <-expired:
	case <-time.After(time.Second):
		t.Fatal("timer did not expire after concurrent resets")
	}

	mu.Lock()
	defer mu.Unlock()
	if expires != 1 {
		t.Fatalf("expiration callbacks = %d, want 1", expires)
	}
}
