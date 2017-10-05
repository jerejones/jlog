package target

import (
	"testing"

	"github.com/jerejones/jlog/event"
)

type testTarget struct {
	writeRx []string
}

func (t *testTarget) Write(entry string) {
	t.writeRx = append(t.writeRx, entry)
}

func TestNewCustomFactory(t *testing.T) {
	tgt := &testTarget{}
	factory := NewCustomFactory(tgt)
	spec := Config{
		Name:       t.Name(),
		Layout:     "${level} ${message}",
		Type:       t.Name(),
		Properties: nil,
	}
	customTgt, err := factory(spec)
	if err != nil {
		t.Errorf("Err: %v", err)
	}
	if customTgt == nil {
		t.Error("Custom Factory returned nil")
	}
	customTgt.Write(event.Info{
		Level:   event.InfoLevel,
		Source:  nil,
		Message: t.Name(),
	})

	expected := "INFO  " + t.Name()
	if expected != tgt.writeRx[0] {
		t.Errorf("Expected %s\nGot %s", expected, tgt.writeRx[0])
	}
}
