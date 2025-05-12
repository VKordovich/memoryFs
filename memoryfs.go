package memoryfs

import (
	"errors"
	"io"
	"io/fs"
	"time"
)

type File struct {
	data    []byte
	name    string
	pos     int64
	modTime time.Time
	closed  bool
}

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (fi *fileInfo) Name() string       { return fi.name }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() fs.FileMode  { return 0644 }
func (fi *fileInfo) ModTime() time.Time { return fi.modTime }
func (fi *fileInfo) IsDir() bool        { return false }
func (fi *fileInfo) Sys() any           { return nil }

var (
	_ fs.File     = &File{}
	_ io.Reader   = &File{}
	_ io.ReaderAt = &File{}
	_ io.Writer   = &File{}
	_ io.WriterAt = &File{}
	_ io.WriterTo = &File{}
	_ io.Seeker   = &File{}
	_ io.Closer   = &File{}
	_ fs.FileInfo = &fileInfo{}
	_ fs.FS       = &FS{}
)

func (f *File) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    f.name,
		size:    int64(len(f.data)),
		modTime: f.modTime,
	}, nil
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.closed {
		return 0, errors.New("file is closed")
	}
	if f.pos >= int64(len(f.data)) {
		return 0, io.EOF
	}
	n = copy(p, f.data[f.pos:])
	f.pos += int64(n)
	return n, nil
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	if f.closed {
		return 0, errors.New("file is closed")
	}
	if off >= int64(len(f.data)) {
		return 0, io.EOF
	}
	n = copy(p, f.data[off:])
	return n, nil
}

func (f *File) Write(p []byte) (n int, err error) {
	if f.closed {
		return 0, errors.New("file is closed")
	}
	f.data = append(f.data, p...)
	return len(p), nil
}

func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	if f.closed {
		return 0, errors.New("file is closed")
	}
	if off+int64(len(p)) > int64(len(f.data)) {
		newData := make([]byte, off+int64(len(p)))
		copy(newData, f.data)
		f.data = newData
	}
	n = copy(f.data[off:], p)
	return n, nil
}

func (f *File) WriteTo(w io.Writer) (n int64, err error) {
	if f.closed {
		return 0, errors.New("file is closed")
	}
	written, err := w.Write(f.data[f.pos:])
	f.pos += int64(written)
	return int64(written), err
}

func (f *File) Seek(off int64, whence int) (int64, error) {
	if f.closed {
		return 0, errors.New("file is closed")
	}
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = off
	case io.SeekCurrent:
		newPos = f.pos + off
	case io.SeekEnd:
		newPos = int64(len(f.data)) + off
	default:
		return 0, errors.New("wrong whence")
	}

	if newPos < 0 {
		return 0, errors.New("wrong position")
	}
	if newPos > int64(len(f.data)) {
		newPos = int64(len(f.data))
	}

	f.pos = newPos
	return newPos, nil
}

func (f *File) Close() error {
	if f.closed {
		return errors.New("file is closed")
	}
	f.closed = true
	return nil
}

func (f *File) Bytes() []byte { return f.data }

type FS struct {
	files map[string][]byte
}

func (fsLocal *FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, errors.New("invalid path")
	}

	if fsLocal.files == nil {
		fsLocal.files = make(map[string][]byte)
	}

	data, exists := fsLocal.files[name]
	if !exists {
		return nil, errors.New("file does not exist")
	}

	return &File{
		data:    data,
		pos:     0,
		name:    name,
		modTime: time.Now(),
		closed:  false,
	}, nil
}
