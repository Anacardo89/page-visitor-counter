package repo

func (r *visitorHandler) AddVisitor(url, visitorID string) {
	r.mu.Lock()
	if _, exists := r.data[url]; !exists {
		r.data[url] = make(map[string]struct{})
	}
	r.data[url][visitorID] = struct{}{}
	r.mu.Unlock()
}

func (r *visitorHandler) CountVisitors(url string) int {
	var visitorCount int
	r.mu.RLock()
	if visitors, exists := r.data[url]; exists {
		visitorCount = len(visitors)
	}
	r.mu.RUnlock()
	return visitorCount
}
