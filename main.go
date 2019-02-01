package main

import (
	"fmt"
	"io"
	"os"
)

// TODO
// - implement flags
// - write mockable test
// - start writing serialised version
// - try concurrency

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
		du:        &serial{},
	}
}

// Run executes the command.
func (c *CLI) Run(args []string) int {
	numFiles, bytes := c.du.Count("dir1", "dir2", "dir3")
	fmt.Fprintf(c.outStream, "%d files, %d bytes", numFiles, bytes)
	return ExitCodeOK
}

func main() {
	cli := New()
	os.Exit(cli.Run(os.Args[1:]))
}

// DiskUsage represents an interface that can count the number of files and total bytes under the directories.
type DiskUsage interface {
	Count(...string) (int, int)
}

type serial struct{}

func (s *serial) Count(dirs ...string) (int, int) {
	return 0, 0
}
