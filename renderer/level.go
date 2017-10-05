package renderer

import (
	"io"

	"github.com/jerejones/jlog/event"
)

func init() {
	RegisterRenderer("level", LayoutFunc(Level))
}

// Level writes info.Level to w
func Level(info event.Info, w io.Writer) (int, error) {
	return w.Write([]byte(info.Level.String()))
}
