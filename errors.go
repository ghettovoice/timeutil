package timeutil

// Error is a package-level error type used for simple sentinel errors.
type Error string

const (
	// ErrInvalidTimerSnapshot is returned when a snapshot is nil or contains
	// invalid timer state and cannot be used to restore a timer.
	ErrInvalidTimerSnapshot Error = "invalid timer snapshot"
)

// Error implements the [error] interface.
func (e Error) Error() string { return string(e) }
