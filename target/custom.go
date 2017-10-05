package target

import (
	"github.com/jerejones/jlog/event"
	"github.com/jerejones/jlog/renderer"
	"github.com/pkg/errors"
)

type custom struct {
	impl  CustomTarget
	implE CustomTargetWithError

	Renderer renderer.Renderer
}

func NewCustomFactory(target CustomTarget) func(spec Config) (Target, error) {
	var implE CustomTargetWithError
	if implError, ok := target.(CustomTargetWithError); ok {
		implE = implError
	}
	return func(spec Config) (Target, error) {
		rndr, err := renderer.New(spec.Layout)
		if err != nil {
			return nil, errors.Wrap(err, "bad layout")
		}

		return &custom{
			Renderer: rndr,

			impl:  target,
			implE: implE,
		}, nil
	}
}

func (t custom) Write(info event.Info) {
	entry := t.Renderer.Render(info)
	if info.Level < event.WarnLevel || t.implE == nil {
		t.impl.Write(entry)
	} else {
		t.implE.WriteError(entry)
	}
}
