// Command icon-builder is a graphical picker for composing app icons —
// pick a background icon, type up to a few letters, see it live. Bare
// invocation opens a local browser UI.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DavidMarsanic/brightencode-appkit/browser"
	"github.com/DavidMarsanic/icon-builder/internal/server"
)

const version = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("icon-builder", flag.ContinueOnError)

	port := fs.Int("port", 0, "local UI server port (default: automatic)")
	showVersion := fs.Bool("version", false, "print the version and exit")
	fs.Usage = func() { printUsage(fs) }

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 1
	}

	if *showVersion {
		fmt.Println("icon-builder " + version)
		return 0
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(ctx)
	addr, err := srv.Start(*port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}

	fmt.Fprintln(os.Stderr, "Icon Builder running at", addr, "— press Ctrl+C to quit")

	// A host process (securexe-launcher) sets SECUREXE_HOSTED before
	// starting us and watches this exact stderr line to discover the URL,
	// so it can host our UI in its own native window instead of a spawned
	// Chrome one. OpenIfNotHosted no-ops in that case.
	if err := browser.OpenIfNotHosted("icon-builder", addr+"/"); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't open a window automatically:", err)
		fmt.Fprintln(os.Stderr, "open this URL manually:", addr+"/")
	}

	<-ctx.Done()
	return 0
}

func printUsage(fs *flag.FlagSet) {
	fmt.Fprint(os.Stderr, `icon-builder — pick a background icon, type up to a few letters, and see
your app icon composed live.

Bare invocation opens the browser UI.

Usage:
  icon-builder              open the browser UI

Flags:
`)
	fs.PrintDefaults()
}
