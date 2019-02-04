package conc

import (
	"io/ioutil"
	"log"
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
	done := make(chan interface{})
	defer close(done)

	resChannels := make([]<-chan result, len(dirs))
	for i, dir := range dirs {
		resChannels[i] = du.walkDir(done, dir)
	}

	num, bytes := 0, 0
	for res := range fanIn(done, resChannels...) {
		if res.err != nil {
			return 0, 0, res.err
		}
		num++
		bytes += res.size
	}

	return num, bytes, nil
}

func (du *DiskUsage) walkDir(done <-chan interface{}, dir string) <-chan result {
	log.Println("pass")
	resCh := make(chan result)

	go func() {
		defer close(resCh)
		// entries, err := du.dirReader(dir)
		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			select {
			case <-done:
			case resCh <- result{size: 0, err: err}:
			}
			return
		}
	outerLoop:
		for _, entry := range entries {
			if entry.IsDir() {
				subdir := filepath.Join(dir, entry.Name())
				for {
					select {
					case <-done:
						return
					case res, ok := <-du.walkDir(done, subdir):
						if !ok {
							continue outerLoop
						}
						select {
						case <-done:
							return
						case resCh <- res:
						}
					}
				}
			} else {
				select {
				case <-done:
					return
				case resCh <- result{size: int(entry.Size()), err: nil}:
				}
			}
		}
	}()

	return resCh
}

func fanIn(
	done <-chan interface{},
	channels ...<-chan result,
) <-chan result {
	var wg sync.WaitGroup
	multiplexedStream := make(chan result)

	multiplex := func(c <-chan result) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
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
