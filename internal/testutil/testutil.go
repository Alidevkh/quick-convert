// Package testutil holds helpers shared by tests: skipping when an engine
// is not installed, and creating small sample files to convert.
package testutil

import (
	"os/exec"
	"testing"
)

// Need skips the test unless every named program is on PATH.
func Need(t *testing.T, programs ...string) {
	t.Helper()
	for _, p := range programs {
		if _, err := exec.LookPath(p); err != nil {
			t.Skipf("%s not installed", p)
		}
	}
}
