package target

import "os"

func init() {
	RegisterTargetFactory("stdout", NewStdOut)
}

type StdOutTarget struct {
	Target
}

func NewStdOut(spec Config) (Target, error) {
	w, err := NewWriter(os.Stdout, spec)
	if err != nil {
		return nil, err
	}
	return &StdOutTarget{Target: w}, nil
}
