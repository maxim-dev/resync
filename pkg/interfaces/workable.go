package interfaces

//go:generate mockgen -source=./workable.go -destination=./workable_mocks.go -package=interfaces

type Workable interface {
	WorkingOn(string)
	WorkingDone(string)
	IsInProgress(string) bool
}
