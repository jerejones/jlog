package target

import (
	"testing"

	"github.com/jerejones/jlog/event"
)

func TestStdErr(t *testing.T) {
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
			tgt, _ := NewStdErr(Config{
				Name:   "stderr",
				Layout: "[${level}] ${message}",
			})

			tgt.Write(event.Info{
				Level:   tc.level,
				Message: "I'm alive!!",
			})

			stdout, stderr := captureStdOutErr(ow, oC, ew, eC)

			assertNotContains(t, stdout, "I'm alive!!")
			assertNotContains(t, stdout, tc.expected)
			assertContains(t, stderr, "I'm alive!!")
			assertContains(t, stderr, tc.expected)

			restoreStdOutErr(origStdOut, origStdErr)
		})
	}
}
