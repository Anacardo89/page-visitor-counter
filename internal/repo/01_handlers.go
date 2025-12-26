package repo

import (
	"sync"
)

type visitorHandler struct {
	mu   sync.RWMutex
	data map[string]map[string]struct{}
}

func NewVisitorRepo() VisitorRepo {
	return &visitorHandler{
		data: make(map[string]map[string]struct{}),
	}
}
