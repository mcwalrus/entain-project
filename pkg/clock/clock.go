// Package clock provides an injectable time interface for testing purposes.
package clock

import (
	"time"
)

var nowFunc func() time.Time = time.Now

// Now returns the current time using the configured time function.
func Now() time.Time {
	return nowFunc()
}

// NowFunc sets a custom function to be used for getting the current time.
func NowFunc(fn func() time.Time) {
	nowFunc = fn
}

// ResetClockImplementation resets the time function back to default time.Now().
func ResetClockImplementation() {
	nowFunc = time.Now
}
