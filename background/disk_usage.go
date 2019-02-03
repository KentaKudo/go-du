package background

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// DiskUsage satisfies an DiskUsage interface.
type DiskUsage struct {
	dirReader func(string) ([]os.FileInfo, error)
}

// New creates a new DiskUsage instance.
func New() *DiskUsage {
	return &DiskUsage{
		dirReader: ioutil.ReadDir,
	}
}

// Count returns the number of files and total bytes under the given directories.
func (du *DiskUsage) Count(dirs []string) (int, int, error) {
	fileSizes := make(chan int)
	go func() {
		for _, d := range dirs {
			du.walkDir(d, fileSizes)
		}
		close(fileSizes)
	}()

	var num, bytes int
	for b := range fileSizes {
		num++
		bytes += b
	}

	return num, bytes, nil
}

func (du *DiskUsage) walkDir(dir string, fileSizes chan<- int) {
	for _, e := range du.dirents(dir) {
		if e.IsDir() {
			subdir := filepath.Join(dir, e.Name())
			du.walkDir(subdir, fileSizes)
		} else {
			fileSizes <- int(e.Size())
		}
	}
}

func (du *DiskUsage) dirents(dir string) []os.FileInfo {
	entries, _ := du.dirReader(dir)
	return entries
}
