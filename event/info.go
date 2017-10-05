package event

type Source interface {
	LoggerName() string
}

type Info struct {
	Message string
	Level   Level
	Source  Source
}
