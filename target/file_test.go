package target

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jerejones/jlog/event"
)

func TestNewFile_MissingFileNameIsError(t *testing.T) {
	f, err := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
	})
	if err == nil {
		t.Error("No error returned for missing filename")
	}
	if f != nil {
		t.Error("Non-nil file target with no filename")
	}
}

func TestNewFile(t *testing.T) {
	f, err := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": "${basedir}/test.log",
		},
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if f == nil {
		t.Error("NewFile returned nil")
	}
}

func TestFile_GetFullFileName(t *testing.T) {
	f, _ := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": "${basedir}/test.log",
		},
	})

	filename := f.(*FileTarget).GetFullFileName(event.Info{})

	if filepath.Base(filename) != "test.log" {
		t.Errorf("Unexpected filename: %s", filename)
	}
}

func TestFile_Write(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testfile_write")
	if err != nil {
		t.Error("Unable to get temp file")
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	fmt.Printf("Filename: %s\n", tmpFile.Name())

	f, _ := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": tmpFile.Name(),
		},
	})

	f.Write(event.Info{
		Level:   event.InfoLevel,
		Message: "Hi Test",
	})

	contents, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Errorf("Unable to read file: %s (%v)", tmpFile.Name(), err)
	}
	expected := "[INFO ] Hi Test\n"
	if string(contents) != expected {
		t.Error("Unexpected contents in temp file:\nWanted: %s\nGot: %s", string(contents), expected)
	}
}

func BenchmarkFile_GetFullFileName_Static(b *testing.B) {
	f, _ := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": "test.log",
		},
	})

	evt := event.Info{
		Level:   event.InfoLevel,
		Message: "Hi Test",
	}

	var res string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		res = f.(*FileTarget).GetFullFileName(evt)
	}
	_ = res
}

func BenchmarkFile_GetFullFileName_Basedir(b *testing.B) {
	f, _ := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": "${basedir}/test.log",
		},
	})

	evt := event.Info{
		Level:   event.InfoLevel,
		Message: "Hi Test",
	}

	var res string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		res = f.(*FileTarget).GetFullFileName(evt)
	}
	_ = res
}

func BenchmarkFile_GetFullFileName_Datetime(b *testing.B) {
	f, _ := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": "${datetime:format=2006-01-02}-test.log",
		},
	})

	evt := event.Info{
		Level:   event.InfoLevel,
		Message: "Hi Test",
	}

	var res string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		res = f.(*FileTarget).GetFullFileName(evt)
	}
	_ = res
}

func BenchmarkFile_GetFullFileName_Level(b *testing.B) {
	f, _ := NewFile(Config{
		Name:   "file",
		Layout: "[${level}] ${message}",
		Properties: map[string]string{
			"filename": "${level}-test.log",
		},
	})

	evt := event.Info{
		Level:   event.InfoLevel,
		Message: "Hi Test",
	}

	var res string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		res = f.(*FileTarget).GetFullFileName(evt)
	}
	_ = res
}
