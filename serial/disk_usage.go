package serial

import (
	"fmt"
	"io/ioutil"
	"os"
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
		files, err := d.dirReader(dir)
		if err != nil {
			return 0, 0, err
		}
		subdirs := []string{}
		for _, f := range files {
			if f.IsDir() {
				subdirs = append(subdirs, fmt.Sprintf("%s/%s", dir, f.Name()))
				continue
			}
			num++
			bytes += int(f.Size())
		}
		subnum, subbytes, err := d.Count(subdirs)
		if err != nil {
			return 0, 0, err
		}
		num += subnum
		bytes += subbytes
	}

	return num, bytes, nil
}
