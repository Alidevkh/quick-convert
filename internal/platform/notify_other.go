//go:build !windows

package platform

import (
	"fmt"
	"os"
)

// Notify prints the message instead of showing a notification (used when
// running outside Windows, e.g. during development).
func Notify(title, body string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", title, body)
}
