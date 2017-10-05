package main

import (
	"time"

	"github.com/jerejones/jlog"
	"github.com/jerejones/jlog/event"
)

func main() {
	l := jlog.GetPackageLogger()
	l.Write(event.InfoLevel, "IT LIVES!")

	for {
		l.Error("Error!!!!")
		l.Warning("Warning")
		l.Info("Info")
		l.Debug("Debug")

		time.Sleep(15 * time.Second)
	}
}
