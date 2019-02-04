package conc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// DiskUsage satisfies DiskUsage interface.
type DiskUsage struct {
	dirReader func(string) ([]os.FileInfo, error)
}

// New creates a new DiskUsage instance.
func New() *DiskUsage {
	return &DiskUsage{dirReader: ioutil.ReadDir}
}

type result struct {
	size int
	err  error
}

// Count counts the number of files and the total bytes under the given directories.
func (du *DiskUsage) Count(dirs []string) (int, int, error) {
	resChannels := make([]<-chan result, len(dirs))
	for i, dir := range dirs {
		resChannels[i] = du.walkDir(dir)
	}

	num, bytes := 0, 0
	for res := range fanIn(resChannels...) {
		if res.err != nil {
			// broadcast
			return 0, 0, res.err
		}
		num++
		bytes += res.size
	}

	return num, bytes, nil
}

func (du *DiskUsage) walkDir(dir string) <-chan result {
	resCh := make(chan result)

	go func() {
		defer close(resCh)
		// entries, err := du.dirReader(dir)
		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			resCh <- result{size: 0, err: err}
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				subdir := filepath.Join(dir, entry.Name())
				for res := range du.walkDir(subdir) {
					resCh <- res
				}
			} else {
				resCh <- result{size: int(entry.Size()), err: nil}
			}
		}
	}()

	return resCh
}

func fanIn(channels ...<-chan result) <-chan result {
	var wg sync.WaitGroup
	multiplexedStream := make(chan result)

	multiplex := func(c <-chan result) {
		defer wg.Done()
		for i := range c {
			multiplexedStream <- i
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
