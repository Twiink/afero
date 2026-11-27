package zipfs

import (
	"archive/zip"
	"io"
	"path/filepath"
	"testing"
)

func TestFileRead(t *testing.T) {
	zrc, err := zip.OpenReader("testdata/small.zip")
	if err != nil {
		t.Fatal(err)
	}
	zfs := New(&zrc.Reader)
	f, err := zfs.Open("smallFile")
	if err != nil {
		t.Fatal(err)
	}
	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	chunkSize := info.Size() * 2 // read with extra large buffer

	buf := make([]byte, chunkSize)
	n, err := f.Read(buf)
	if err != io.EOF {
		t.Fatal("Failed to read file to completion:", err)
	}
	if n != int(info.Size()) {
		t.Errorf("Expected read length to be %d, found: %d", info.Size(), n)
	}

	// read a second time to check f.offset and f.buf are correct
	buf = make([]byte, chunkSize)
	n, err = f.Read(buf)
	if err != io.EOF {
		t.Fatal("Failed to read a fully read file:", err)
	}
	if n != 0 {
		t.Errorf("Expected read length to be 0, found: %d", n)
	}
}

func TestFileNameAfterClose(t *testing.T) {
	archive, err := zip.OpenReader("testdata/t.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	fs := New(&archive.Reader)
	for _, path := range []string{"testFile", "testDir1/testFile", "testDir1", "sub", "/"} {
		t.Run(path, func(t *testing.T) {
			file, err := fs.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			name := filepath.Join(string(filepath.Separator), filepath.FromSlash(path))
			if got := file.Name(); got != name {
				t.Fatalf("Name() = %q, want %q", got, name)
			}
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}
			if got := file.Name(); got != name {
				t.Errorf("Name() after Close() = %q, want %q", got, name)
			}
		})
	}
}
