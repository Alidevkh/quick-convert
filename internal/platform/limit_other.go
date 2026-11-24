//go:build !windows

package platform

// AcquireHeavy does not limit anything outside Windows.
func AcquireHeavy() (release func()) {
	return func() {}
}
