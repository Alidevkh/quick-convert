// genscripts writes install.cmd, uninstall.cmd and selftest.cmd into a folder.
// Used by scripts/build.sh:  go run ./cmd/genscripts dist/QuickConvert
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"quickconvert/internal/menu"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: genscripts <folder>")
		os.Exit(2)
	}
	dir := os.Args[1]
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fail(err)
	}
	files := map[string]string{
		"install.cmd":   menu.InstallScript(),
		"uninstall.cmd": menu.UninstallScript(),
		"selftest.cmd":  menu.SelfTestScript(),
	}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
			fail(err)
		}
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
