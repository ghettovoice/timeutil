// Package timeutil provides serializable timer types that behave like
// [time.AfterFunc] but can be snapshotted, serialized, and restored across
// process restarts or long-lived workflows.
//
// The main type is [Timer], which keeps its runtime behaviour,
// including automatic callback execution via a background [time.Timer], while
// exposing deterministic state through [TimerSnapshot].
// Snapshots can be marshaled to JSON, stored, and later passed to [RestoreTimer]
// to obtain a fresh timer instance. Callbacks are runtime-only; they must be
// reattached after restoration with [Timer.SetCallback] or
// [Timer.Reset].
//
// [Watchdog] is a complementary helper that fires a callback after a
// configured period of inactivity unless it is reset.
//
// Basic usage:
//
//	// Create timer with automatic callback execution.
//	timer := timeutil.AfterFunc(5*time.Second, func() {
//	    log.Println("Timer expired!")
//	})
//
//	// Persist the timer state.
//	snap := timer.Snapshot()
//	data, _ := json.Marshal(snap)
//
//	// Restore later and reattach callbacks.
//	var restoredSnap timeutil.TimerSnapshot
//	_ = json.Unmarshal(data, &restoredSnap)
//	restored, _ := timeutil.RestoreTimer(&restoredSnap)
//	restored.SetCallback(func() {
//	    log.Println("Restored timer expired!")
//	})
//
// All timer operations are thread-safe and can be called concurrently from
// multiple goroutines.
package timeutil
