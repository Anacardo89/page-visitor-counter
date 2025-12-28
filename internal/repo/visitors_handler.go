package repo

import (
	"sync"
)

type visitorRepoHandler struct {
	mu   sync.RWMutex
	data map[string]map[string]struct{}
}

func NewVisitorRepo() VisitorRepo {
	return &visitorRepoHandler{
		data: make(map[string]map[string]struct{}),
	}
}

func (r *visitorRepoHandler) AddVisitor(url, visitorID string) {
	r.mu.Lock()
	if _, exists := r.data[url]; !exists {
		r.data[url] = make(map[string]struct{})
	}
	r.data[url][visitorID] = struct{}{}
	r.mu.Unlock()
}

func (r *visitorRepoHandler) CountVisitors(url string) int {
	var visitorCount int
	r.mu.RLock()
	if visitors, exists := r.data[url]; exists {
		visitorCount = len(visitors)
	}
	r.mu.RUnlock()
	return visitorCount
}
