//go:build !windows

package platform

// Stand-in for the heavy-job limit when running outside Windows.

// AcquireHeavy does not limit anything outside Windows.
func AcquireHeavy() (release func()) {
	return func() {}
}
