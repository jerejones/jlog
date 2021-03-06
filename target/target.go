package target

import (
	"fmt"

	"github.com/jerejones/jlog/event"
	"github.com/pkg/errors"
)

var targetFactories map[string]Factory

type Factory func(Config) (Target, error)

type CustomTarget interface {
	Write(string)
}

type CustomTargetWithError interface {
	CustomTarget
	WriteError(string)
}

type Target interface {
	Write(info event.Info)
}

type UnknownTargetError string

func (err UnknownTargetError) Error() string {
	return fmt.Sprintf("Unknown target type: %s", err)
}

func RegisterTargetFactory(targetType string, factory Factory) error {
	if targetFactories == nil {
		targetFactories = make(map[string]Factory)
	}

	_, exists := targetFactories[targetType]
	if exists {
		return errors.New("Type already registered")
	}
	targetFactories[targetType] = factory
	return nil

}

func New(spec Config) (Target, error) {
	if _, exists := targetFactories[spec.Type]; !exists {
		return nil, UnknownTargetError(spec.Type)
	}
	target, err := targetFactories[spec.Type](spec)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating target %s", spec.Name)
	}

	return target, nil
}

func RegisterCustomTarget(targetType string, target CustomTarget) error {
	return RegisterTargetFactory(targetType, NewCustomFactory(target))
}
