package dispenser

import "sync"

type WorkingList struct {
	storage map[string]bool
	mu      sync.Mutex
}

func NewWorkingList() *WorkingList {
	return &WorkingList{
		storage: make(map[string]bool),
		mu:      sync.Mutex{},
	}
}

func (r *WorkingList) WorkingOn(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[path] = true
}

func (r *WorkingList) WorkingDone(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, path)
}

func (r *WorkingList) IsInProgress(path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.storage[path]

	return exists
}