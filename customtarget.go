package jlog

import "github.com/jerejones/jlog/target"

// CustomTarget is the interface that every custom target has to implement.
//
// Write will be called with the rendered log entry
// Note that string will not end with a line feed
type CustomTarget interface {
	Write(string)
}

// CustomTargetWithError is the interface that custom targets may implement so that WARN and ERROR level
// log entries may be handled differently
//
// Note that string will not end with a line feed
type CustomTargetWithError interface {
	CustomTarget
	WriteError(string)
}

// RegisterCustomLoggerTarget allows a custom target to be used in your config
//
// name is the value that will go in the Type field of the target config
func RegisterCustomLoggerTarget(name string, logger CustomTarget) error {
	target.RegisterCustomTarget(name, logger)
	return nil
}
