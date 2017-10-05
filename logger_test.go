package jlog_test

import (
	"reflect"
	"testing"

	"github.com/jerejones/jlog"
)

func TestGetNamedLogger(t *testing.T) {
	l := jlog.GetNamedLogger("test")
	if l.LoggerName() != "test" {
		t.Errorf(`Expected "test" as logger name. Got %s`, l.LoggerName())
	}
}

func TestGetPackageLogger(t *testing.T) {
	l := jlog.GetPackageLogger()
	if l.LoggerName() != "github.com/jerejones/jlog" {
		t.Errorf(`Expected "github.com/jerejones/jlog" as logger name. Got %s`, l.LoggerName())
	}
}

func TestGetObjectLogger(t *testing.T) {
	l := jlog.GetObjectLogger(t)
	if l.LoggerName() != "testing.T" {
		t.Errorf(`Expected "testing.T" as logger name. Got %s`, l.LoggerName())
	}
}

func TestGetTypeLogger(t *testing.T) {
	l := jlog.GetTypeLogger(reflect.TypeOf(0))
	if l.LoggerName() != "int" {
		t.Errorf(`Expected "int" as logger name. Got %s`, l.LoggerName())
	}
}

func TestLogger_WriteEvent(t *testing.T) {
	/*
		b := &bytes.Buffer{}
		l := jlog.GetNamedLogger("WriteEventTest")
		w, err := target.NewWriter(b, target.Config{
			Name:   "WriteEventTest",
			Layout: "${datetime} [${level}] ${message}",
		})
		assert.Nil(t, err)
			l.Targets = append(l.Targets, w)
			l.WriteEvent(event.Info{
				Level:   event.InfoLevel,
				Message: "It LIVES!",
			})

			logged := b.String()
			assert.NotZero(t, len(logged))
			fmt.Println(logged)
	*/
}
