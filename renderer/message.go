package renderer

import (
	"io"

	"github.com/jerejones/jlog/event"
)

func init() {
	RegisterRenderer("message", LayoutFunc(Message))

}

// Message writes info.Message to w
func Message(info event.Info, w io.Writer) (int, error) {
	return w.Write([]byte(info.Message))
}
