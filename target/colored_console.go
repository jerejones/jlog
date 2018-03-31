package target

import (
	"os"

	"github.com/jerejones/jlog/event"
)

func init() {
	RegisterTargetFactory("coloredconsole", NewColoredConsole)
}

type ColoredConsoleTarget struct {
	stderr Target
	stdout Target
}

func NewColoredConsole(spec Config) (Target, error) {
	stderr, err := NewColorWriter(os.Stderr, spec)
	if err != nil {
		return nil, err
	}
	stdout, err := NewColorWriter(os.Stdout, spec)
	if err != nil {
		return nil, err
	}

	return &ColoredConsoleTarget{
		stderr: stderr,
		stdout: stdout,
	}, nil
}

func (t *ColoredConsoleTarget) Write(info event.Info) {
	if info.Level < event.ErrorLevel {
		t.stdout.Write(info)
	} else {
		t.stderr.Write(info)
	}
}
