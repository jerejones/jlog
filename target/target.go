package target

import (
	"github.com/jerejones/jlog/event"
	"github.com/pkg/errors"
)

var targetFactories map[string]Factory

type Factory func(Config) (Target, error)

type Target interface {
	Write(info event.Info)
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
		return nil, errors.Errorf("Unknown target type: %s", spec.Type)
	}
	target, err := targetFactories[spec.Type](spec)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating target %s", spec.Name)
	}

	return target, nil
}
