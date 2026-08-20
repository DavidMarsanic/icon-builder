// Package browser opens URLs in the OS-native way, so the app can hand
// off to the user's browser the same way whether it's opening its own UI
// in app-mode or a plain link.
package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open launches the user's default browser at target.
func Open(target string) error {
	return run(openArgs(target))
}

func openArgs(target string) []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"open", target}
	case "windows":
		return []string{"rundll32", "url.dll,FileProtocolHandler", target}
	default:
		return []string{"xdg-open", target}
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command to run")
	}
	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("running %s: %w", args[0], err)
	}
	go cmd.Wait() // reap the process without blocking the caller
	return nil
}
