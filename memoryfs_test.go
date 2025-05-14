package memoryfs

import (
	"errors"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestFileRead(t *testing.T) {
	f := &File{
		data:    []byte("Want Avanpost"),
		pos:     0,
		name:    "test.txt",
		modTime: time.Now(),
		closed:  false,
	}
	type testCase struct {
		name     string
		bufSize  int
		wantRead string
		wantN    int
		wantErr  error
		wantPos  int64
	}
	tests := []testCase{
		{
			name:     "read content",
			bufSize:  15,
			wantRead: "Want Avanpost",
			wantN:    13,
			wantErr:  nil,
			wantPos:  13,
		},
		{
			name:     "read at EOF",
			bufSize:  10,
			wantRead: "",
			wantN:    0,
			wantErr:  io.EOF,
			wantPos:  13,
		},
		{
			name:     "read after close",
			bufSize:  5,
			wantRead: "",
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			wantPos:  0,
		},
	}

	for i, tt := range tests {
		if i == 2 {
			f.Close()
		}
		t.Run(tt.name, func(t *testing.T) {
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

// Без теста закрытого файла
func TestFileReadAt(t *testing.T) {
	f := &File{
		data:    []byte("Want Avanpost"),
		pos:     0,
		name:    "test.txt",
		modTime: time.Now(),
		closed:  false,
	}

	type testCase struct {
		name     string
		off      int64
		bufSize  int
		wantRead string
		wantN    int
		wantErr  error
	}

	tests := []testCase{
		{
			name:     "read from start",
			off:      0,
			bufSize:  4,
			wantRead: "Want",
			wantN:    4,
			wantErr:  nil,
		},
		{
			name:     "read from middle",
			off:      5,
			bufSize:  8,
			wantRead: "Avanpost",
			wantN:    8,
			wantErr:  nil,
		},
		{
			name:     "read at EOF",
			off:      13,
			bufSize:  5,
			wantRead: "",
			wantN:    0,
			wantErr:  io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	f := &File{
		data:    []byte("Avanpost"),
		pos:     0,
		name:    "test.txt",
		modTime: time.Now(),
		closed:  false,
	}

	type testCase struct {
		name     string
		data     []byte
		wantN    int
		wantErr  error
		wantData string
	}

	tests := []testCase{
		{
			name:     "write data",
			data:     []byte(" Moscow"),
			wantN:    7,
			wantErr:  nil,
			wantData: "Avanpost Moscow",
		},
		{
			name:     "write after close",
			data:     []byte("Ignored text"),
			wantN:    0,
			wantErr:  errors.New("file is closed"),
			wantData: "Avanpost Moscow",
		},
	}

	for i, tt := range tests {
		if i == 1 {
			f.Close()
		}
		t.Run(tt.name, func(t *testing.T) {
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

// Без теста закрытого файла
func TestFileWriteAt(t *testing.T) {
	f := &File{
		data:    []byte("Avanpost Moscow"),
		pos:     0,
		name:    "test.txt",
		modTime: time.Now(),
		closed:  false,
	}
	type testCase struct {
		name     string
		data     []byte
		off      int64
		wantN    int
		wantErr  error
		wantData string
	}

	tests := []testCase{
		{
			name:     "write at end",
			data:     []byte("London"),
			off:      9,
			wantN:    6,
			wantErr:  nil,
			wantData: "Avanpost London",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
