package renderer

import (
	"io"

	"github.com/jerejones/jlog/event"
)

func init() {
	RegisterRenderer("logger_name", LayoutFunc(LoggerName))
}

// LoggerName writes info.Source.LoggerName() to w
func LoggerName(info event.Info, w io.Writer) (int, error) {
	return w.Write([]byte(info.Source.LoggerName()))
}
