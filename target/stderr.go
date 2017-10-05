package target

import "os"

func init() {
	RegisterTargetFactory("stderr", NewStdErr)
}

type StdErrTarget struct {
	Target
}

func NewStdErr(spec Config) (Target, error) {
	w, err := NewWriter(os.Stderr, spec)
	if err != nil {
		return nil, err
	}
	return &StdErrTarget{Target: w}, nil
}
