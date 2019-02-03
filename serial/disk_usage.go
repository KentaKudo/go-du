package serial

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// DiskUsage represents an instance that satisfies DiskUsage interface.
type DiskUsage struct {
	dirReader func(string) ([]os.FileInfo, error)
}

// New creates a new DiskUsage instance.
func New() *DiskUsage {
	return &DiskUsage{dirReader: ioutil.ReadDir}
}

// Count counts the number of files and sizes under the given directories with serial access.
func (d *DiskUsage) Count(dirs []string) (int, int, error) {
	var num, bytes int
	for _, dir := range dirs {
		n, b, err := d.walkDir(dir)
		if err != nil {
			return 0, 0, err
		}
		num += n
		bytes += b
	}

	return num, bytes, nil
}

func (d *DiskUsage) walkDir(dir string) (int, int, error) {
	var num, bytes int
	entries, err := d.dirReader(dir)
	if err != nil {
		return 0, 0, err
	}
	for _, e := range entries {
		if e.IsDir() {
			n, b, err := d.walkDir(filepath.Join(dir, e.Name()))
			if err != nil {
				return 0, 0, nil
			}
			num += n
			bytes += b
			continue
		}

		num++
		bytes += int(e.Size())
	}

	return num, bytes, nil
}
