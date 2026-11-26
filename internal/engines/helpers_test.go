package engines

import "os"

func writeExecutable(path string) error {
	return os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755)
}
