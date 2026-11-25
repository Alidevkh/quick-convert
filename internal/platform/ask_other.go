//go:build !windows

package platform

// Stand-in for the "Which pages?" box outside Windows.

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AskPages asks on the terminal (used when running outside Windows).
func AskPages(pages int) (answer string, ok bool, err error) {
	fmt.Fprint(os.Stderr, AskPrompt(pages)+"\n> ")
	line, readErr := bufio.NewReader(os.Stdin).ReadString('\n')
	answer = strings.TrimSpace(line)
	return answer, readErr == nil && answer != "", nil
}
