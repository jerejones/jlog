package jlog

import (
	"reflect"
	"testing"

	"github.com/jerejones/jlog/event"
	"github.com/jerejones/jlog/target"
)

func TestNewRouter(t *testing.T) {
	spec := RouterConfig{
		Name:     "*",
		MinLevel: "info",
		MaxLevel: "warn",
	}
	r := newRouter(spec)

	expected := &router{
		pattern: "*",
		levels: map[event.Level]bool{
			event.DebugLevel: false,
			event.InfoLevel:  true,
			event.WarnLevel:  true,
			event.ErrorLevel: false,
		},
		targets: []target.Target{},
	}

	if !reflect.DeepEqual(expected, r) {
		t.Error("Error creating new router")
	}
}
