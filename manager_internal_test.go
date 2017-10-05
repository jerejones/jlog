package jlog

import (
	"testing"
)

func TestManager_updateLogger(t *testing.T) {
	m, err := NewManager(DefaultConfig())
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	l := &logger{name: "Test"}
	m.updateLogger(l)

	if len(l.routers) != 1 {
		t.Errorf("Unexpected number of routers. Got %d, expected %d", len(l.routers), 1)
	}
}
