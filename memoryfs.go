package memoryfs

import (
	"io"
	"io/fs"
)

type File struct {
	data []byte
}

var (
	_ fs.File     = &File{}
	_ io.Reader   = &File{}
	_ io.ReaderAt = &File{}
	_ io.Writer   = &File{}
	_ io.WriterAt = &File{}
	_ io.WriterTo = &File{}
	_ io.Seeker   = &File{}
	_ io.Closer   = &File{}
)

func (f *File) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Seek(off int64, whence int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) WriteTo(w io.Writer) (n int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Close() error {
	//TODO implement me
	panic("implement me")
}

func (f *File) Stat() (fs.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Bytes() []byte { return f.data }

type FS struct {
}

var _ fs.FS = &FS{}

func (fs *FS) Open(name string) (fs.File, error) {
	//TODO implement me
	panic("implement me")
}
