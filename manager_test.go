package jlog_test

import (
	"testing"

	"github.com/jerejones/jlog"
)

func TestNewManager(t *testing.T) {
	m, err := jlog.NewManager(jlog.DefaultConfig())
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if m == nil {
		t.Error("Unexpected nil")
	}
}

func TestManager_ApplyConfig(t *testing.T) {
	m, err := jlog.NewManager(nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	err = m.ApplyConfig(jlog.DefaultConfig())
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
