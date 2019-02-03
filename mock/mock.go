package mock

var DefaultDiskUsage = &DiskUsage{CountFn: func(_ []string) (int, int, error) { return 0, 0, nil }}

type DiskUsage struct {
	CountFn      func([]string) (int, int, error)
	CountInvoked bool
}

func (d *DiskUsage) Count(dirs []string) (int, int, error) {
	d.CountInvoked = true
	return d.CountFn(dirs)
}
