package target

import (
	"encoding/json"
	"strings"
)

type Config struct {
	Name       string
	Layout     string
	Type       string
	Properties map[string]string
}

func (tc *Config) UnmarshalJSON(b []byte) error {
	var props map[string]string
	err := json.Unmarshal(b, &props)
	if err != nil {
		return err
	}

	tc.Properties = make(map[string]string)
	for key, val := range props {
		switch strings.ToLower(key) {
		case "name":
			tc.Name = val
		case "type":
			tc.Type = val
		case "layout":
			tc.Layout = val
		default:
			tc.Properties[key] = val
		}
	}

	return nil
}
