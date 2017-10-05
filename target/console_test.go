package target

import (
	"testing"

	"github.com/jerejones/jlog/event"
)

func TestConsole(t *testing.T) {
	tt := []struct {
		name           string
		level          event.Level
		expectedStdOut string
		expectedStdErr string
	}{
		{name: "debug to stdout", level: event.DebugLevel, expectedStdOut: "DEBUG"},
		{name: "info to stdout", level: event.InfoLevel, expectedStdOut: "INFO "},
		{name: "warn to stdout", level: event.WarnLevel, expectedStdOut: "WARN "},
		{name: "error to stderr", level: event.ErrorLevel, expectedStdErr: "ERROR"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			origStdOut, origStdErr, ow, oC, ew, eC := startStdOutErrCapture()
			tgt, _ := NewConsole(Config{
				Name:   "stdout",
				Layout: "[${level}] ${message}",
			})

			tgt.Write(event.Info{
				Level:   tc.level,
				Message: "I'm alive!!",
			})

			stdout, stderr := captureStdOutErr(ow, oC, ew, eC)

			if len(tc.expectedStdErr) > 0 {
				assertNotContains(t, stdout, "I'm alive!!")
				assertNotContains(t, stdout, tc.expectedStdErr)
				assertContains(t, stderr, "I'm alive!!")
				assertContains(t, stderr, tc.expectedStdErr)
			} else {
				assertContains(t, stdout, "I'm alive!!")
				assertContains(t, stdout, tc.expectedStdOut)
				assertNotContains(t, stderr, "I'm alive!!")
				assertNotContains(t, stderr, tc.expectedStdOut)
			}
			restoreStdOutErr(origStdOut, origStdErr)
		})
	}
}
