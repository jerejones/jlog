package jlog

import (
	"runtime"
	"strings"
)

// CurrentPackageFileName returns the name of the current file in the current package
//
// For example calling it in this file would return "github.com/jerejones/jlog/namers.go"
func CurrentPackageFileName() string {
	pc, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	start := strings.Index(filename, packageName)
	if start < 0 {
		return "unknown"
	}
	return filename[start:]
}
