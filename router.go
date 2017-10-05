package jlog

import (
	"fmt"
	"strings"

	"github.com/jerejones/jlog/event"
	"github.com/jerejones/jlog/target"
)

type router struct {
	pattern string
	levels  map[event.Level]bool
	targets []target.Target
}

func newRouter(spec RouterConfig) *router {
	levels := make(map[event.Level]bool)
	minLevel := event.NewLevel(spec.MinLevel)
	if minLevel == event.UnknownLevel {
		minLevel = event.DebugLevel
	}
	maxLevel := event.NewLevel(spec.MaxLevel)
	if maxLevel == event.UnknownLevel {
		maxLevel = event.ErrorLevel
	}
	for _, level := range event.AllLevels() {
		levels[level] = (level >= minLevel) && (level <= maxLevel)
	}
	return &router{
		pattern: spec.Name,
		levels:  levels,
		targets: []target.Target{},
	}
}

// Matches returns whether the router's pattern matches the name provided
//
// The code was adapted from http://github.com/ryanuber/go-glob/network
// which is MIT licensed: https://github.com/ryanuber/go-glob/blob/master/LICENSE
func (r *router) Matches(name string) bool {
	// Empty pattern can only match empty subject
	if r.pattern == "" {
		return name == r.pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if r.pattern == "*" {
		return true
	}

	parts := strings.Split(r.pattern, "*")

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return name == r.pattern
	}

	leadingGlob := strings.HasPrefix(r.pattern, "*")
	trailingGlob := strings.HasSuffix(r.pattern, "*")
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(name, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		name = name[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(name, parts[end])
}

func (r *router) Write(source event.Source, level event.Level, v ...interface{}) {
	if enabled, exists := r.levels[level]; enabled && exists {
		info := event.Info{
			Level:   level,
			Source:  source,
			Message: fmt.Sprint(v...),
		}
		for _, t := range r.targets {
			t.Write(info)
		}
	}
}

func (r *router) Writef(source event.Source, level event.Level, format string, a ...interface{}) {
	if enabled, exists := r.levels[level]; enabled && exists {
		info := event.Info{
			Level:   level,
			Source:  source,
			Message: fmt.Sprintf(format, a...),
		}
		for _, t := range r.targets {
			t.Write(info)
		}
	}
}
