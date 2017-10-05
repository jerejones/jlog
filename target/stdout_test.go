package target

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/jerejones/jlog/event"
)

func TestStdOut(t *testing.T) {
	tt := []struct {
		name     string
		level    event.Level
		expected string
	}{
		{name: "debug", level: event.DebugLevel, expected: "[DEBUG]"},
		{name: "info", level: event.InfoLevel, expected: "[INFO ]"},
		{name: "warn", level: event.WarnLevel, expected: "[WARN ]"},
		{name: "error", level: event.ErrorLevel, expected: "[ERROR]"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			origStdOut, origStdErr, ow, oC, ew, eC := startStdOutErrCapture()
			tgt, _ := NewStdOut(Config{
				Name:   "stdout",
				Layout: "[${level}] ${message}",
			})

			tgt.Write(event.Info{
				Level:   tc.level,
				Message: "I'm alive!!",
			})

			stdout, stderr := captureStdOutErr(ow, oC, ew, eC)

			assertContains(t, stdout, "I'm alive!!")
			assertContains(t, stdout, tc.expected)
			assertNotContains(t, stderr, "I'm alive!!")
			assertNotContains(t, stderr, tc.expected)

			restoreStdOutErr(origStdOut, origStdErr)
		})
	}
}

func startStdOutErrCapture() (*os.File, *os.File, *os.File, chan string, *os.File, chan string) {
	origStdOut := os.Stdout
	origStdErr := os.Stderr
	or, ow, _ := os.Pipe()
	os.Stdout = ow
	oC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, or)
		oC <- buf.String()
	}()
	er, ew, _ := os.Pipe()
	os.Stderr = ew
	eC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, er)
		eC <- buf.String()
	}()
	return origStdOut, origStdErr, ow, oC, ew, eC
}

func captureStdOutErr(ow *os.File, oC chan string, ew *os.File, eC chan string) (string, string) {
	ow.Close()
	stdout := <-oC
	ew.Close()
	stderr := <-eC
	return stdout, stderr
}

func restoreStdOutErr(origStdOut *os.File, origStdErr *os.File) {
	os.Stdout = origStdOut
	os.Stderr = origStdErr
}
