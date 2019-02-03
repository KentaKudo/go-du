package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/KentaKudo/go-du/mock"
)

func TestExitCode(t *testing.T) {
	sut := &CLI{
		outStream: new(bytes.Buffer),
		errStream: new(bytes.Buffer),
		du:        mock.DefaultDiskUsage,
	}
	input := strings.Split("test", " ")
	want := ExitCodeOK
	if got := sut.Run(input); got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestRun_DiskUsageInvokation(t *testing.T) {
	mock := mock.DefaultDiskUsage
	sut := &CLI{
		outStream: new(bytes.Buffer),
		errStream: new(bytes.Buffer),
		du:        mock,
	}
	input := strings.Split("test1 test2", " ")
	sut.Run(input)
	if !mock.CountInvoked {
		t.Errorf("CLI.du.Count is not called")
	}
}
