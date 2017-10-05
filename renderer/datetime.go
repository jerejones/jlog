package renderer

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/jerejones/jlog/event"
)

const (
	defaultDateTimeFormat = "2006-01-02 15:04:05.999"
)

var (
	_ LayoutRenderer               = (*DateTime)(nil)
	_ LayoutRendererWithParameters = (*DateTime)(nil)
)

func init() {
	RegisterRenderer("datetime", &DateTime{format: defaultDateTimeFormat})
}

type DateTime struct {
	format string
	utc    bool
}

func (dt *DateTime) Write(info event.Info, w io.Writer) (int, error) {
	t := time.Now()
	if dt.utc {
		t = t.UTC()
	}
	return w.Write([]byte(time.Now().Format(dt.format)))
}

func (*DateTime) IsDynamic() bool {
	return true
}

func (dt *DateTime) CreateLayoutRenderer(parameters map[string]string) (LayoutRenderer, error) {
	ret := DateTime{
		format: defaultDateTimeFormat,
	}
	for key, val := range parameters {
		switch strings.ToLower(key) {
		case "format":
			ret.format = val
		case "utc":
			if len(val) == 0 {
				ret.utc = true
			} else {
				utc, err := strconv.ParseBool(val)
				if err != nil {
					return nil, fmt.Errorf("invalid value for utc: %s", val)
				}
				ret.utc = utc
			}

		default:
			return nil, fmt.Errorf("unknown paramater: %s", key)
		}
	}
	return &ret, nil
}
