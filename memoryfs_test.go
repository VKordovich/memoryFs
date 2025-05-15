package memoryfs

import (
	"errors"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestFileRead(t *testing.T) {
	type testCase struct {
		name     string
		bufSize  int
		wantRead string
		wantN    int
		wantErr  error
		wantPos  int64
		closed   bool
	}

	tests := []testCase{
		{
			name:     "read content",
			bufSize:  15,
			wantRead: "Want Avanpost",
			wantN:    13,
			wantErr:  nil,
			wantPos:  13,
			closed:   false,
		},
		{
			name:     "read at EOF",
			bufSize:  10,
			wantRead: "",
			wantN:    0,
			wantErr:  io.EOF,
			wantPos:  13,
			closed:   false,
		},
		{
			name:     "read after close",
			bufSize:  5,
			wantRead: "",
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			wantPos:  0,
			closed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				data:    []byte("Want Avanpost"),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}

			if tt.closed {
				f.Close()
			}

			if tt.name == "read at EOF" {
				f.pos = int64(len(f.data))
			}

			buf := make([]byte, tt.bufSize)
			n, err := f.Read(buf)

			if n != tt.wantN {
				t.Errorf("TestFileRead > n = %d, want %d", n, tt.wantN)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFileRead > err = %v, want %v", err, tt.wantErr)
			}
			if strRes := string(buf[:n]); strRes != tt.wantRead {
				t.Errorf("TestFileRead > content = %q, want %q", strRes, tt.wantRead)
			}
			if f.pos != tt.wantPos {
				t.Errorf("File.pos = %d, want %d", f.pos, tt.wantPos)
			}
		})
	}
}

func TestFileReadAt(t *testing.T) {
	type testCase struct {
		name     string
		off      int64
		bufSize  int
		wantRead string
		wantN    int
		wantErr  error
		closed   bool
	}

	tests := []testCase{
		{
			name:     "read from start",
			off:      0,
			bufSize:  4,
			wantRead: "Want",
			wantN:    4,
			wantErr:  nil,
			closed:   false,
		},
		{
			name:     "read from middle",
			off:      5,
			bufSize:  8,
			wantRead: "Avanpost",
			wantN:    8,
			wantErr:  nil,
			closed:   false,
		},
		{
			name:     "read at EOF",
			off:      13,
			bufSize:  5,
			wantRead: "",
			wantN:    0,
			wantErr:  io.EOF,
			closed:   false,
		},
		{
			name:     "read after close",
			off:      0,
			bufSize:  0,
			wantRead: "",
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			closed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				data:    []byte("Want Avanpost"),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}
			if tt.closed {
				f.Close()
			}

			buf := make([]byte, tt.bufSize)
			n, err := f.ReadAt(buf, tt.off)
			if n != tt.wantN {
				t.Errorf("TestFileReadAt n = %d, want %d", n, tt.wantN)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFileReadAt err = %v, want %v", err, tt.wantErr)
			}
			if strRes := string(buf[:n]); strRes != tt.wantRead {
				t.Errorf("TestFileReadAt content = %q, want %q", strRes, tt.wantRead)
			}
		})
	}
}

func TestFileWrite(t *testing.T) {
	type testCase struct {
		name     string
		data     []byte
		wantN    int
		wantErr  error
		wantData string
		closed   bool
	}

	tests := []testCase{
		{
			name:     "write data",
			data:     []byte(" Moscow"),
			wantN:    7,
			wantErr:  nil,
			wantData: "Avanpost Moscow",
			closed:   false,
		},
		{
			name:     "write after close",
			data:     []byte("Ignored text"),
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			wantData: "Avanpost",
			closed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				data:    []byte("Avanpost"),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}

			if tt.closed {
				f.Close()
			}
			n, err := f.Write(tt.data)
			if n != tt.wantN {
				t.Errorf("TestFileWrite n = %d, want %d", n, tt.wantN)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFileWrite err = %v, want %v", err, tt.wantErr)
			}
			if strRes := string(f.data); strRes != tt.wantData {
				t.Errorf("File.data = %q, want %q", strRes, tt.wantData)
			}
		})
	}
}

func TestFileWriteAt(t *testing.T) {
	type testCase struct {
		name     string
		data     []byte
		off      int64
		wantN    int
		wantErr  error
		wantData string
		closed   bool
	}

	tests := []testCase{
		{
			name:     "write at end",
			data:     []byte("London"),
			off:      9,
			wantN:    6,
			wantErr:  nil,
			wantData: "Avanpost London",
			closed:   false,
		},
		{
			name:     "write after close",
			data:     []byte("Ignored text"),
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			wantData: "Avanpost Moscow",
			closed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				data:    []byte("Avanpost Moscow"),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}
			if tt.closed {
				f.Close()
			}
			n, err := f.WriteAt(tt.data, tt.off)
			if n != tt.wantN {
				t.Errorf("TestFileWriteAt n = %d, want %d", n, tt.wantN)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFileWriteAt err = %v, want %v", err, tt.wantErr)
			}
			if strRes := string(f.data); strRes != tt.wantData {
				t.Errorf("File.data = %q, want %q", strRes, tt.wantData)
			}
		})
	}
}

func TestFileWriteTo(t *testing.T) {
	type testCase struct {
		name     string
		wantN    int64
		wantErr  error
		wantData string
		closed   bool
	}

	tests := []testCase{
		{
			name:     "write content",
			wantN:    15,
			wantErr:  nil,
			wantData: "Avanpost Moscow",
			closed:   false,
		},
		{
			name:     "write after close",
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			wantData: "",
			closed:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				data:    []byte("Avanpost Moscow"),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}

			fw := &File{
				data:    []byte(""),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}

			if tt.closed {
				f.Close()
			}

			n, err := f.WriteTo(fw)
			if n != tt.wantN {
				t.Errorf("TestFileWriteTo n = %d, want %d", n, tt.wantN)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFileWriteTo err = %v, want %v", err, tt.wantErr)
			}
			if strRes := string(fw.data); strRes != tt.wantData {
				t.Errorf("TestFileWriteTo data = %q, want %q", strRes, tt.wantData)
			}
		})
	}
}

func TestFileSeek(t *testing.T) {
	type testCase struct {
		name    string
		off     int64
		whence  int
		wantPos int64
		wantErr error
		closed  bool
	}

	tests := []testCase{
		{
			name:    "seek to start",
			off:     0,
			whence:  io.SeekStart,
			wantPos: 0,
			wantErr: nil,
			closed:  false,
		},
		{
			name:    "seek to middle",
			off:     6,
			whence:  io.SeekStart,
			wantPos: 6,
			wantErr: nil,
			closed:  false,
		},
		{
			name:    "seek relative current",
			off:     2,
			whence:  io.SeekCurrent,
			wantPos: 2,
			wantErr: nil,
			closed:  false,
		},
		{
			name:    "seek to end",
			off:     0,
			whence:  io.SeekEnd,
			wantPos: 15,
			wantErr: nil,
			closed:  false,
		},
		{
			name:    "seek beyond end",
			off:     16,
			whence:  io.SeekStart,
			wantPos: 15,
			wantErr: nil,
			closed:  false,
		},
		{
			name:    "seek negative",
			off:     -1,
			whence:  io.SeekStart,
			wantPos: 0,
			wantErr: errors.New("wrong position"),
			closed:  false,
		},
		{
			name:    "invalid whence",
			off:     0,
			whence:  99,
			wantPos: 0,
			wantErr: errors.New("wrong whence"),
			closed:  false,
		},
		{
			name:    "seek after close",
			off:     0,
			whence:  99,
			wantPos: 0,
			wantErr: errors.New("file is closed"),
			closed:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				data:    []byte("Avanpost Moscow"),
				pos:     0,
				name:    "test.txt",
				modTime: time.Now(),
				closed:  false,
			}

			if tt.closed {
				f.Close()
			}

			newPos, err := f.Seek(tt.off, tt.whence)
			if newPos != tt.wantPos {
				t.Errorf("TestFileSeek newPos = %d, want %d", newPos, tt.wantPos)
			}
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFileSeek err = %v, want %v", err, tt.wantErr)
			}
			if f.pos != tt.wantPos {
				t.Errorf("File.pos = %d, want %d", f.pos, tt.wantPos)
			}
		})
	}
}

func TestFileClose(t *testing.T) {
	f := &File{
		data:    []byte("Avanpost Moscow"),
		pos:     0,
		name:    "test.txt",
		modTime: time.Now(),
		closed:  false,
	}

	err := f.Close()
	if err != nil {
		t.Errorf("TestFileClose error on first call: %v", err)
	}
	if !f.closed {
		t.Errorf("File.closed = false, want true after Close")
	}
	if f.pos != 0 {
		t.Errorf("File.pos != 0, want 0 after Close")
	}

	err = f.Close()
	if !reflect.DeepEqual(err, errors.New("file is closed")) {
		t.Errorf("TestFileClose error on second call = %v, want %v", err, errors.New("file is closed"))
	}
}

func TestFileStat(t *testing.T) {
	f := &File{
		data:    []byte("Hello"),
		name:    "test.txt",
		modTime: time.Now(),
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}

	if name := info.Name(); name != "test.txt" {
		t.Errorf("Name() = %q, want %q", name, "test.txt")
	}
	if size := info.Size(); size != 5 {
		t.Errorf("Size() = %d, want %d", size, 5)
	}
}

func TestFSOpen(t *testing.T) {
	type testCase struct {
		name    string
		path    string
		wantErr error
	}

	tests := []testCase{
		{
			name:    "valid file",
			path:    "test.txt",
			wantErr: nil,
		},
		{
			name:    "non-existent file",
			path:    "missing.txt",
			wantErr: errors.New("file does not exist"),
		},
		{
			name:    "invalid path",
			path:    "../test.txt",
			wantErr: errors.New("invalid path"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsLocal := &FS{
				files: map[string][]byte{
					"test.txt": []byte(""),
				},
			}
			_, err := fsLocal.Open(tt.path)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TestFSOpen err = %v, want %v", err, tt.wantErr)
				return
			}
		})
	}
}
