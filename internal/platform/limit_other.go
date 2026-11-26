//go:build !windows

package platform

// Stand-in for the job limit when running outside Windows.

import "quickconvert/internal/actions"

// Acquire does not limit anything outside Windows.
func Acquire(actions.Kind) (release func()) {
	return func() {}
}
