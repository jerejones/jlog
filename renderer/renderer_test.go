package renderer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jerejones/jlog/event"
)

func TestLevelRenderer(t *testing.T) {
	for _, level := range []event.Level{event.ErrorLevel, event.InfoLevel, event.DebugLevel} {
		info := event.Info{
			Level: level,
		}
		buf := new(bytes.Buffer)
		Level(info, buf)

		expected := level.String()
		output := buf.String()

		if output != expected {
			t.Errorf("Expected: %s Got: %s", expected, output)
		}
	}
}

type mockSource struct {
	name string
}

func (s mockSource) LoggerName() string {
	return s.name
}

func TestLoggerNameRenderer(t *testing.T) {
	info := event.Info{
		Source: mockSource{name: "*"},
	}
	buf := new(bytes.Buffer)
	LoggerName(info, buf)
	output := buf.String()
	if output != "*" {
		t.Errorf("Expected: * Got: %s", output)
	}
}

func TestMessageRenderer(t *testing.T) {
	msg := "Test 1 2 3"
	info := event.Info{
		Message: msg,
	}
	buf := new(bytes.Buffer)
	Message(info, buf)
	output := buf.String()
	if output != msg {
		t.Errorf("Expected: %s Got: %s", msg, output)
	}
}

func TestDateTimeRenderer_NoParameters(t *testing.T) {
	r, err := New("${datetime}")
	if err != nil {
		t.Error("Unable to create renderer")
	}
	info := event.Info{}
	output := r.Render(info)

	fmt.Println(output)
}

func TestDateTimeRenderer_FormatParam(t *testing.T) {
	r, err := New("${datetime:format=2006-01-02}")
	if err != nil {
		t.Error("Unable to create renderer")
	}
	info := event.Info{}
	output := r.Render(info)

	fmt.Println(output)
}
func TestDateTimeRenderer_UtcParam(t *testing.T) {

	r, err := New("${datetime:utc=0} ${datetime:utc=1} ${datetime:utc}")
	if err != nil {
		t.Error("Unable to create renderer")
	}
	info := event.Info{}
	output := r.Render(info)

	fmt.Println(output)
}
func TestDateTimeRenderer_UnknownParam(t *testing.T) {
	_, err := New("${datetime:bob=0}")
	if err == nil {
		t.Error("Unexpectedly able to create renderer")
	}

}

func TestSplitTag(t *testing.T) {
	tag, params := splitTag("datetime")
	if tag != "datetime" {
		t.Errorf("Wrong tag. Expected \"datetime\", got \"%s\"", tag)
	}
	if params != nil {
		t.Errorf("Wrong params")
	}

	tag, params = splitTag("datetime:format=\"2006-'01'-02\",utc=true,test='\"日本語\"'")
	if tag != "datetime" {
		t.Errorf("Wrong tag. Expected \"datetime\", got \"%s\"", tag)
	}
	if params == nil {
		t.Errorf("Wrong params")
	}

}
