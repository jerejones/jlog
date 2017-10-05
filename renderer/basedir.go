package renderer

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jerejones/jlog/event"
)

var (
	_ LayoutRenderer = (*BaseDir)(nil)
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	bd := BaseDir{
		dir: []byte(exPath),
	}
	RegisterRenderer("basedir", &bd)
}

// BaseDir renders the directory of the executable relative to the current working
// directory and is typically useful for specifying where to save log files
//
// Example:
//     "filename": "${basedir}/logs/debug.log"
type BaseDir struct {
	dir []byte
}

// Write writes the directory of the executable to w
func (bd *BaseDir) Write(info event.Info, w io.Writer) (int, error) {
	return w.Write(bd.dir)
}

// IsDynamic returns false
func (*BaseDir) IsDynamic() bool {
	return false
}
