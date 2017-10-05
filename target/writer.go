package target

import (
	"io"

	"github.com/jerejones/jlog/event"
	"github.com/jerejones/jlog/renderer"
	"github.com/pkg/errors"
)

type WriterTarget struct {
	io.Writer
	Layout   string
	Renderer renderer.Renderer
}

func NewWriter(w io.Writer, spec Config) (Target, error) {
	rndr, err := renderer.New(spec.Layout)
	if err != nil {
		return nil, errors.Wrap(err, "bad layout")

	}
	return &WriterTarget{
		Writer:   w,
		Layout:   spec.Layout,
		Renderer: rndr,
	}, nil
}

func (t *WriterTarget) Write(info event.Info) {
	t.Writer.Write([]byte(t.Renderer.Render(info) + "\n"))
}
