package target

import (
	"os"

	"github.com/jerejones/jlog/event"
)

func init() {
	RegisterTargetFactory("console", NewConsole)
}

type ConsoleTarget struct {
	stderr Target
	stdout Target
}

func NewConsole(spec Config) (Target, error) {
	stderr, err := NewWriter(os.Stderr, spec)
	if err != nil {
		return nil, err
	}
	stdout, err := NewWriter(os.Stdout, spec)
	if err != nil {
		return nil, err
	}
	return &ConsoleTarget{
		stderr: stderr,
		stdout: stdout,
	}, nil
}

func (t *ConsoleTarget) Write(info event.Info) {
	if info.Level < event.ErrorLevel {
		t.stdout.Write(info)
	} else {
		t.stderr.Write(info)
	}
}
