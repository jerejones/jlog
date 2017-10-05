package jlog

import (
	"github.com/jerejones/jlog/event"
)

var (
	_ Logger = (*logger)(nil)
)

type Logger interface {
	Error(v ...interface{}) error
	Warning(v ...interface{}) error
	Info(v ...interface{}) error
	Debug(v ...interface{}) error

	Errorf(format string, a ...interface{}) error
	Warningf(format string, a ...interface{}) error
	Infof(format string, a ...interface{}) error
	Debugf(format string, a ...interface{}) error
}

type logger struct {
	name    string
	routers []*router
}

func (l *logger) LoggerName() string {
	return l.name
}

func (l *logger) Write(level event.Level, v ...interface{}) {
	for _, r := range l.routers {
		r.Write(l, level, v...)
	}
}

func (l *logger) Writef(level event.Level, format string, a ...interface{}) {
	for _, r := range l.routers {
		r.Writef(l, level, format, a...)
	}
}

func (l *logger) Error(v ...interface{}) error {
	l.Write(event.ErrorLevel, v...)
	return nil
}

func (l *logger) Warning(v ...interface{}) error {
	l.Write(event.WarnLevel, v...)
	return nil
}

func (l *logger) Info(v ...interface{}) error {
	l.Write(event.InfoLevel, v...)
	return nil
}

func (l *logger) Debug(v ...interface{}) error {
	l.Write(event.DebugLevel, v...)
	return nil
}

func (l *logger) Errorf(format string, a ...interface{}) error {
	l.Writef(event.ErrorLevel, format, a...)
	return nil
}

func (l *logger) Warningf(format string, a ...interface{}) error {
	l.Writef(event.WarnLevel, format, a...)
	return nil
}

func (l *logger) Infof(format string, a ...interface{}) error {
	l.Writef(event.InfoLevel, format, a...)
	return nil
}

func (l *logger) Debugf(format string, a ...interface{}) error {
	l.Writef(event.DebugLevel, format, a...)
	return nil
}
