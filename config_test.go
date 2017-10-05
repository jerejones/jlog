package jlog_test

import (
	"strings"
	"testing"

	"github.com/jerejones/jlog"
)

func assertStringsEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Error("Strings not equal\nExpected: %s\nGot: %s", expected, actual)
	}
}

func assertIntsEqual(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Error("Ints not equal\nExpected: %d\nGot: %d", expected, actual)
	}
}

func TestLoadConfig_Json(t *testing.T) {
	input := `{
	"autoreload": true,
	"routes": [
		{
			"name": "*",
			"targets": "c,f",
			"minlevel": "info"
		}
	],
	"targets": [
		{
			"name": "c",
			"type": "console",
			"layout": "${date} {$message}"
		}, {
			"name": "f",
			"type": "file",
			"layout": "${date} [${level}] {$message}",
			"filename": "${basedir}/test.log",
			"header": ">>>start log",
			"max_open_files": "5"
		}
	]
}`

	cfg, err := jlog.UnmarshalConfig(strings.NewReader(input))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if !cfg.AutoReload {
		t.Error("Error: Autoload is false")
	}
	assertIntsEqual(t, 1, len(cfg.Routes))
	assertIntsEqual(t, 2, len(cfg.Targets))

	target1 := cfg.Targets[0]
	assertStringsEqual(t, "c", target1.Name)
	assertStringsEqual(t, "console", target1.Type)
	assertStringsEqual(t, "${date} {$message}", target1.Layout)

	target2 := cfg.Targets[1]
	assertStringsEqual(t, "f", target2.Name)
	assertStringsEqual(t, "file", target2.Type)
	assertStringsEqual(t, "${date} [${level}] {$message}", target2.Layout)
	assertStringsEqual(t, "${basedir}/test.log", target2.Properties["filename"])
	assertStringsEqual(t, ">>>start log", target2.Properties["header"])
	assertStringsEqual(t, "5", target2.Properties["max_open_files"])
}
