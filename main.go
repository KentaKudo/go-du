package main

import (
	"fmt"
	"os"
)

// TODO
// - apply command best practice
// - implement flags
// - write mockable test
// - start writing serialised version
// - try concurrency

// DiskUsage represents an interface that can count the number of files and total bytes under the directories.
type DiskUsage interface {
	Count(dir ...string) (int, int)
}

type serial struct{}

func (s *serial) Count(dir ...string) (int, int) {
	return 0, 0
}

func main() {
	s := &serial{}
	numFiles, bytes := s.Count("dir1", "dir2", "dir3")
	fmt.Fprintf(os.Stdout, "%d files, %d bytes", numFiles, bytes)
}
