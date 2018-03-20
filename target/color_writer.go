package target

import (
	"io"

	"github.com/fatih/color"
	"github.com/jerejones/jlog/event"
	"github.com/jerejones/jlog/renderer"
	"github.com/pkg/errors"
)

type ColorWriterTarget struct {
	io.Writer
	Layout   string
	Renderer renderer.Renderer
	colors   map[event.Level]*color.Color
}

func NewColorWriter(w io.Writer, spec Config) (Target, error) {
	rndr, err := renderer.New(spec.Layout)
	if err != nil {
		return nil, errors.Wrap(err, "bad layout")
	}

	colors := map[event.Level]*color.Color{
		event.ErrorLevel:   color.New(color.FgRed),
		event.DebugLevel:   color.New(color.FgHiBlue),
		event.WarnLevel:    color.New(color.FgYellow),
		event.InfoLevel:    color.New(color.FgWhite),
		event.UnknownLevel: color.New(color.FgWhite),
	}

	return &ColorWriterTarget{
		Writer:   w,
		Layout:   spec.Layout,
		Renderer: rndr,
		colors:   colors,
	}, nil
}

func (t *ColorWriterTarget) Write(info event.Info) {
	if c, ok := t.colors[info.Level]; ok {
		c.Fprintf(t.Writer, t.Renderer.Render(info)+"\n")
		return
	}
	t.colors[event.UnknownLevel].Fprintf(t.Writer, t.Renderer.Render(info)+"\n")
}
