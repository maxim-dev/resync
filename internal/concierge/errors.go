package concierge

import "fmt"

type EmptyPath struct{}

func (e *EmptyPath) Error() string {
	return "a path is empty string"
}

type IsNotADir struct {
	name string
}

func (e *IsNotADir) Error() string {
	return fmt.Sprintf("%q is not a directory", e.name)
}

func (e *IsNotADir) Is(target error) bool {
	te, ok := target.(*IsNotADir)
	if ok == false {
		return false
	}
	return e.name == te.name
}

type TheSamePath struct{}

func (e *TheSamePath) Error() string {
	return "src and dest are the same path"
}
