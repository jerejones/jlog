package event

import (
	"strings"
)

// Level indicates the severity of the log event
type Level int

const (
	UnknownLevel Level = 1 << iota
	// DebugLevel is used for debugging information.  Executed queries, user authenticated, session expired
	DebugLevel
	// InfoLevel is used for informative purposes that are not errors. Normal behavior like mail sent, user updated profile etc.
	InfoLevel
	// WarnLevel is used when something went wrong but the application can continue
	WarnLevel
	// ErrorLevel is used when something went wrong and needs attention.
	ErrorLevel
)

func AllLevels() []Level {
	return []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}
}

// String returns level as a string that is 5 characters long
func (level Level) String() string {
	switch level {
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN "
	case InfoLevel:
		return "INFO "
	case DebugLevel:
		return "DEBUG"
	}
	return "UNKNOWN"
}

func NewLevel(str string) Level {
	str = strings.TrimSpace(str)
	if str == "" {
		return UnknownLevel
	}

	switch strings.ToUpper(str[0:1])[0] {
	case 'E':
		return ErrorLevel
	case 'W':
		return WarnLevel
	case 'I':
		return InfoLevel
	case 'D':
		return DebugLevel
	}
	return UnknownLevel
}
