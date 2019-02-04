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
	resCh := make(chan result)
	var n sync.WaitGroup
	for _, dir := range dirs {
		n.Add(1)
		go du.walkDir(dir, &n, resCh)
	}
	go func() {
		n.Wait()
		close(resCh)
	}()

	num, bytes := 0, 0
	for res := range resCh {
		if res.err != nil {
			return 0, 0, res.err
		}
		num++
		bytes += res.size
	}

	return num, bytes, nil
}

func (du *DiskUsage) walkDir(dir string, n *sync.WaitGroup, resCh chan<- result) {
	defer n.Done()
	entries, err := du.dirents(dir)
	if err != nil {
		resCh <- result{size: 0, err: err}
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go du.walkDir(subdir, n, resCh)
		} else {
			resCh <- result{size: int(entry.Size()), err: nil}
		}
	}
}

func (du *DiskUsage) dirents(dir string) ([]os.FileInfo, error) {
	return du.dirReader(dir)
}
