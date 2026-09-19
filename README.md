# timeutil

[![Go Reference](https://pkg.go.dev/badge/github.com/ghettovoice/timeutil.svg)](https://pkg.go.dev/github.com/ghettovoice/timeutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/ghettovoice/timeutil)](https://goreportcard.com/report/github.com/ghettovoice/timeutil)
[![Tests](https://github.com/ghettovoice/timeutil/actions/workflows/test.yml/badge.svg)](https://github.com/ghettovoice/timeutil/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/ghettovoice/timeutil/badge.svg?branch=master)](https://coveralls.io/github/ghettovoice/timeutil?branch=master)
[![CodeQL](https://github.com/ghettovoice/timeutil/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/ghettovoice/timeutil/actions/workflows/github-code-scanning/codeql)

`timeutil` is a small Go package that provides serializable timers with automatic
callback execution. The timers behave like the standard [time.AfterFunc](https://pkg.go.dev/time#AfterFunc),
but their deterministic state (start time, duration, and current state) can be
snapshotted and serialized for persistence, inspection, or transfer between
processes.

## Features

- **Snapshot state** — marshal and unmarshal timer state to/from JSON.
- **Automatic callback execution** — real `time.Timer` runs in the background; no
  manual polling is required.
- **State-aware** — tracks `running`, `stopped`, and `expired` states.
- **AfterFunc support** — create timers that execute a callback when they expire.
- **Thread-safe** — all public methods are safe for concurrent use.
- **Snapshot support** — export lightweight `TimerSnapshot` values
  for persistence outside JSON.
- **Watchdog timer** — a separate helper that fires after a period of inactivity
  unless reset.

## Installation

```bash
go get github.com/ghettovoice/timeutil
```

Requires Go 1.25 or later.

## Usage

### Basic timer

```go
package main

import (
    "fmt"
    "time"

    "github.com/ghettovoice/timeutil"
)

func main() {
    // Create a new timer that expires after 5 seconds.
    timer := timeutil.NewTimer(5 * time.Second)

    // Check whether the timer has expired.
    if timer.Expired() {
        fmt.Println("Timer has expired")
    }

    // Time remaining.
    fmt.Println("Time remaining:", timer.Left())

    // Time elapsed.
    fmt.Println("Time elapsed:", timer.Elapsed())
}
```

### Timer with a callback

```go
timer := timeutil.AfterFunc(5*time.Second, func() {
    fmt.Println("Timer expired!")
})

// Or set a callback after creation.
timer := timeutil.NewTimer(5 * time.Second)
timer.SetCallback(func() {
    fmt.Println("Timer expired!")
})
```

`SetCallback` automatically executes the callback immediately if the timer has
already expired, and it starts a real `time.Timer` for timers that are still
running.

### Serialization

```go
// Create a timer.
timer := timeutil.NewTimer(10 * time.Second)

// Serialize to JSON.
data, err := timer.ToJSON()
if err != nil {
    panic(err)
}

// Restore later.
restored, err := timeutil.FromJSON(data)
if err != nil {
    panic(err)
}

// Callbacks are not serialized; reattach them after restoration.
restored.SetCallback(func() {
    fmt.Println("Restored timer expired!")
})
```

The serialized representation is the same as `TimerSnapshot`:

```json
{
  "start_time": "2025-11-03T12:54:05.184256+03:00",
  "duration": 5000000000,
  "state": "running",
  "stop_time": "2025-11-03T12:54:07.184256+03:00"
}
```

- `start_time` — ISO 8601 timestamp when the timer started.
- `duration` — duration in nanoseconds.
- `state` — one of `running`, `stopped`, or `expired`.
- `stop_time` — present only for stopped timers.

### Recreating a timer from a specific start time

Use `FromTime` when you need to recreate a timer from saved timing metadata:

```go
startTime := time.Now().Add(-2 * time.Second)
timer := timeutil.FromTime(startTime, 5*time.Second)

// UpdateState checks whether the timer has already expired and triggers the
// callback if one is set.
timer.UpdateState()
```

### Watchdog timer

```go
watchdog := timeutil.NewWatchdog(30*time.Second, func() {
    fmt.Println("No activity for 30 seconds")
})

// Later, to keep the watchdog from firing:
watchdog.Reset()

// To stop it permanently:
watchdog.Stop()
```

## API overview

### Timer

Creation:

- `NewTimer(duration)` — create and start a timer immediately.
- `AfterFunc(duration, callback)` — create a timer that executes a callback on expiration.
- `FromTime(startTime, duration)` — recreate a timer with a specific start time.
- `FromJSON(data)` — deserialize a timer from JSON.

State queries:

- `State()` — current timer state.
- `StartTime()` — when the timer started.
- `Duration()` — configured duration.
- `StopTime()` — when the timer was stopped (zero value if not stopped).
- `Expired()` — whether the timer has expired.
- `Elapsed()` — elapsed time since start.
- `Left()` — remaining time until expiration.

Control:

- `Stop()` — stop the timer and prevent callback execution.
- `Reset(duration)` — restart the timer with a new duration, preserving any callback.
- `SetCallback(callback)` — attach or replace the expiration callback.
- `UpdateState()` — re-check expiration and trigger callbacks if needed.

Serialization:

- `Snapshot()` — capture the current state as a `TimerSnapshot`.
- `SnapshotTimer(timer)` — same as `timer.Snapshot()`, but safe on `nil`.
- `RestoreTimer(snapshot)` — create a new timer from a snapshot.
- `ToJSON()` — serialize to JSON.
- Implements `json.Marshaler` and `json.Unmarshaler`.

### Watchdog

- `NewWatchdog(timeout, callback)` — create a watchdog that fires after `timeout`.
- `Reset()` — postpone firing until `timeout` from now.
- `Stop()` — stop the watchdog permanently.

## Important notes

1. **Automatic execution** — when a callback is set, a real `time.Timer` runs in
   the background, so `UpdateState()` is usually not needed.
2. **`UpdateState()` is called automatically during JSON unmarshaling**, so timers
   restored with `FromJSON` already reflect the current time.
3. **`SetCallback()` checks for expiration automatically** — it runs the callback
   immediately if the timer is already expired.
4. **Call `UpdateState()` manually** only when you need to re-check expiration
   after time has passed, or after creating a timer with `FromTime` before a
   callback was set.
5. Callbacks are executed in their own goroutine, like `time.AfterFunc`.
6. Callbacks are **not serialized** and must be reattached after restoration.
7. `Left()` returns `0` for stopped or expired timers.
8. `Reset()` preserves any existing callback but restarts the real timer with the
   new duration.

## Thread safety

All public methods on `Timer` and `Watchdog` are safe for
concurrent use from multiple goroutines without external synchronization.

## Contributing

Contributions are welcome. Please make sure all checks pass before submitting a
pull request:

```bash
task check
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file
for details.
