package main

import (
	"fmt"
	"io"
	"os"

	"github.com/KentaKudo/go-du/background"
)

// TODO
// - try concurrency
// - benchmarking

const (
	// ExitCodeOK is returned when the command runs successfully.
	ExitCodeOK int = iota
	// ExitCodeError is returned when an error occurs.
	ExitCodeError
)

// CLI represents a command line interface which holds dependencies.
type CLI struct {
	outStream, errStream io.Writer

	du DiskUsage
}

// New returns a new CLI instance.
func New() *CLI {
	return &CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
		du:        background.New(),
	}
}

// Run executes the command.
func (c *CLI) Run(args []string) int {
	num, bytes, err := c.du.Count(args)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	fmt.Fprintf(c.outStream, "%d files, %d bytes\n", num, bytes)
	return ExitCodeOK
}

func main() {
	cli := New()
	os.Exit(cli.Run(os.Args[1:]))
}

// DiskUsage represents an interface that can count the number of files and total bytes under the directories.
type DiskUsage interface {
	Count([]string) (int, int, error)
}
