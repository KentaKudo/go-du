package mock

type DiskUsage struct {
	CountFn      func(...string) (int, int)
	CountInvoked bool
}

func (d *DiskUsage) Count(dirs ...string) (int, int) {
	d.CountInvoked = true
	return d.CountFn(dirs...)
}
