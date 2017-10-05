package target

import (
	"bytes"
	"testing"

	"github.com/jerejones/jlog/event"
)

func TestWriter(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	w, err := NewWriter(buf, Config{
		Layout: "${level} ${message}\n",
	})
	if err != nil {
		t.Errorf("Error return from NewWriter: %s", err)
	}

	str := buf.String()

	info := event.Info{
		Message: "Test Message",
		Level:   event.ErrorLevel,
	}

	w.Write(info)

	str = buf.String()
	assertContains(t, str, "Test Message")
}
