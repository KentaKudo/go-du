package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KentaKudo/go-du/background"
	"github.com/KentaKudo/go-du/concurrent"
	"github.com/KentaKudo/go-du/serial"
)

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
func New(strategy string) *CLI {
	var du DiskUsage
	switch strategy {
	case "serial":
		du = serial.New()
	case "background":
		du = background.New()
	case "concurrent":
		du = concurrent.New()
	default:
		panic("unknown strategy")
	}

	return &CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
		du:        du,
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
	var strategy string
	flag.StringVar(&strategy, "strategy", "serial", "The way to access directories")
	flag.Parse()

	cli := New(strategy)
	os.Exit(cli.Run(flag.Args()))
}

// DiskUsage represents an interface that can count the number of files and total bytes under the directories.
type DiskUsage interface {
	Count([]string) (int, int, error)
}
