package renderer

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/jerejones/jlog/event"
	"github.com/pkg/errors"
	"github.com/valyala/fasttemplate"
)

const (
	startTag string = "${"
	endTag   string = "}"
)

var (
	renderers map[string]LayoutRenderer
)

// LayoutRenderer is the interface that wraps a layout renderer.
//
// When creating custom renderers, this is the interface that needs to be implemented.
type LayoutRenderer interface {
	// Write takes the information in info and writes it to w.  The return values are
	// the number of bytes written and any error that occurred.
	Write(info event.Info, w io.Writer) (int, error)
	// IsDynamic returns false if the data written by Write is static throughout the
	// course of the applications lifetime. It returns true otherwise.
	IsDynamic() bool
}

// LayoutRendererWithParameters is the interface which wraps a LayoutRenderer factory
// such that parameters can be specified in the layout
//
// Custom LayoutRenderers will typically also implement this interface if its behaviour
// can be customized
type LayoutRendererWithParameters interface {
	CreateLayoutRenderer(parameters map[string]string) (LayoutRenderer, error)
}

// RegisterRenderer takes a name and an instance of a LayoutRenderer and registers it so
// that the name can be used in layout templates.
//
// If r also implements LayoutRendererWithParameters then r.CreateLayoutRenderer will be
// called to create a new instance of the LayoutRenderer if the template specifies any
// parameters
func RegisterRenderer(name string, r LayoutRenderer) {
	if renderers == nil {
		renderers = make(map[string]LayoutRenderer)
	}
	renderers[name] = r
}

// The LayoutFunc type is an adapter to allow the use of ordinary functions as layout
// renderers. If f is a function with the appropriate signature, LayoutFunc(f) is a
// LayoutRenderer that calls f.
type LayoutFunc func(event.Info, io.Writer) (int, error)

// Write calls f(info, w)
func (f LayoutFunc) Write(info event.Info, w io.Writer) (int, error) {
	return f(info, w)
}

// IsDynamic returns true
func (f LayoutFunc) IsDynamic() bool {
	return true
}

// Renderer is the interface that wraps a specific template
type Renderer interface {
	// Render returns a string rendered from info based on the template provided to New
	Render(info event.Info) string
}

// New returns a Renderer that conforms to the template provided using registered
// LayoutRenderers and an error if there was an error in the template
func New(template string) (Renderer, error) {
	tags, err := extractTags(template)
	if err != nil {
		return nil, err
	}

	isDynamic := false
	for _, tag := range tags {
		r, err := createRenderer(tag)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to create renderer from template \"%s\"", template))
		}
		if r.IsDynamic() {
			isDynamic = true
		}
	}

	ft, err := fasttemplate.NewTemplate(template, startTag, endTag)
	if err != nil {
		return nil, err
	}

	dr := dynamicRenderer{
		template: ft,
	}

	if isDynamic {
		return &dr, nil
	}

	return &staticRenderer{
		text: dr.Render(event.Info{}),
	}, nil
}

func createRenderer(fullTag string) (LayoutRenderer, error) {
	tag, params := splitTag(fullTag)
	baseRenderer, exists := renderers[tag]
	if !exists {
		return nil, fmt.Errorf("unrecognized layout %s", tag)
	}
	if params == nil {
		return baseRenderer, nil
	}
	creator, ok := baseRenderer.(LayoutRendererWithParameters)
	if !ok {
		return nil, fmt.Errorf("%s does not take parameters", tag)
	}
	r, err := creator.CreateLayoutRenderer(params)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to create layout for %s", tag))
	}
	renderers[fullTag] = r
	return r, nil
}

func extractTags(template string) ([]string, error) {
	s := []byte(template)
	a := []byte(startTag)
	b := []byte(endTag)
	tags := []string{}

	tagsCount := bytes.Count(s, a)
	if tagsCount == 0 {
		return nil, nil
	}

	if tagsCount > cap(tags) {
		tags = make([]string, 0, tagsCount)
	}

	for {
		n := bytes.Index(s, a)
		if n < 0 {
			break
		}

		s = s[n+len(a):]
		n = bytes.Index(s, b)
		if n < 0 {
			return nil, fmt.Errorf("cannot find end tag=%q in the template=%q starting from %q", endTag, template, s)
		}

		tags = append(tags, string(s[:n]))
		s = s[n+len(b):]
	}
	return tags, nil
}

func splitTag(fullTag string) (string, map[string]string) {
	colonPos := strings.Index(fullTag, ":")
	if colonPos == -1 {
		return fullTag, nil
	}
	tag := fullTag[:colonPos]
	parameters := make(map[string]string)

	paramName := ""
	paramVal := ""
	currQuotes := rune(0)
	inName := true

	for _, r := range fullTag[colonPos+1:] {
		c := rune(r)
		inQuotes := currQuotes != 0
		if c == currQuotes {
			currQuotes = rune(0)
			continue
		}
		if !inQuotes && (c == '"' || c == '\'') {
			currQuotes = c
			continue
		}
		if !inQuotes && c == ',' {
			if len(paramName) > 0 {
				parameters[paramName] = paramVal
			}
			paramName = ""
			paramVal = ""
			inName = true
			continue
		}
		if !inQuotes && c == '=' {
			inName = false
			continue
		}
		if inName {
			paramName += string(c)
		} else {
			paramVal += string(c)
		}

	}
	if len(paramName) > 0 {
		parameters[paramName] = paramVal
	}

	return tag, parameters
}

type dynamicRenderer struct {
	template *fasttemplate.Template
}

func (r *dynamicRenderer) Render(info event.Info) string {
	return r.template.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		renderer, exists := renderers[tag]
		if !exists {
			return 0, fmt.Errorf("Unknown tag: %s", tag)
		}
		return renderer.Write(info, w)
	})
}

type staticRenderer struct {
	text string
}

func (r *staticRenderer) Render(info event.Info) string {
	return r.text
}
