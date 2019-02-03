package background

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDiskUsage_CountEmpty(t *testing.T) {
	sut := New()
	want := 0
	num, size, err := sut.Count([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if num != want {
		t.Errorf("want %d, got %d", want, num)
	}
	if size != want {
		t.Errorf("want %d, got %d", want, size)
	}
}

// func TestDiskUsage_CountDirReadError(t *testing.T) {
// 	want := "test error"
// 	mock := func(got string) ([]os.FileInfo, error) {
// 		return nil, errors.New(want)
// 	}
// 	sut := &DiskUsage{dirReader: mock}
// 	if _, _, err := sut.Count([]string{"dir1"}); err.Error() != want {
// 		t.Errorf("want %q, got %q", want, err.Error())
// 	}
// }

func TestDiskUsage_CountEmptyDir(t *testing.T) {
	sut := New()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	want := 0
	num, size, err := sut.Count([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if num != want {
		t.Errorf("want %d, got %d", want, num)
	}
	if size != want {
		t.Errorf("want %d, got %d", want, size)
	}
}

func TestDiskUsage_CountOneFile(t *testing.T) {
	sut := New()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	f, err := ioutil.TempFile(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	sizeWant, err := f.Write([]byte("0"))
	if err != nil {
		t.Fatal(err)
	}

	numWant := 1
	num, size, err := sut.Count([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if num != numWant {
		t.Errorf("want %d, got %d", numWant, num)
	}
	if size != sizeWant {
		t.Errorf("want %d, got %d", sizeWant, size)
	}
}

func TestDiskUsage_CountOneDirWithOneFile(t *testing.T) {
	sut := New()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	subdir, err := ioutil.TempDir(dir, "subdir")
	if err != nil {
		t.Fatal(err)
	}
	f, err := ioutil.TempFile(subdir, "test")
	if err != nil {
		t.Fatal(err)
	}
	sizeWant, err := f.Write([]byte("0"))
	if err != nil {
		t.Fatal(err)
	}

	numWant := 1
	num, size, err := sut.Count([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if num != numWant {
		t.Errorf("want %d, got %d", numWant, num)
	}
	if size != sizeWant {
		t.Errorf("want %d, got %d", sizeWant, size)
	}
}
