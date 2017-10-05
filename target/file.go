package target

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/jerejones/jlog/event"
	"github.com/jerejones/jlog/renderer"
	"github.com/pkg/errors"
)

func init() {
	RegisterTargetFactory("file", NewFile)
}

type FileTarget struct {
	Header           string
	Layout           string
	MsgRenderer      renderer.Renderer
	FileNameRenderer renderer.Renderer
}

func NewFile(spec Config) (Target, error) {
	if len(spec.Properties["filename"]) == 0 {
		return nil, errors.New("bad filename")
	}
	msgRndr, err := renderer.New(spec.Layout)
	if err != nil {
		return nil, errors.Wrap(err, "bad layout")
	}
	fnRndr, err := renderer.New(spec.Properties["filename"])
	if err != nil {
		return nil, errors.Wrap(err, "bad filename")
	}
	lruSize, err := strconv.Atoi(spec.Properties["max_open_files"])
	if err != nil || lruSize == 0 {
		lruSize = 5
	}
	f := &FileTarget{
		Header:           spec.Properties["header"],
		Layout:           spec.Layout,
		MsgRenderer:      msgRndr,
		FileNameRenderer: fnRndr,
	}

	return f, nil
}

func (t *FileTarget) Write(info event.Info) {
	filename := t.GetFullFileName(info)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	f.Write([]byte(t.MsgRenderer.Render(info) + "\n"))

	f.Close()
}

func (t *FileTarget) GetFullFileName(info event.Info) string {
	return filepath.Clean(t.FileNameRenderer.Render(info))
}
