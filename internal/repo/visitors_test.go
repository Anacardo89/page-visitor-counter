package repo

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitorHandler_AddVisitor(t *testing.T) {
	tests := []struct {
		name         string
		operations   [][2]string
		expectedData map[string][]string
	}{
		{
			name: "single visitor",
			operations: [][2]string{
				{"https://example.com", "user1"},
			},
			expectedData: map[string][]string{
				"https://example.com": {"user1"},
			},
		},
		{
			name: "duplicate visitor",
			operations: [][2]string{
				{"https://example.com", "user1"},
				{"https://example.com", "user1"},
			},
			expectedData: map[string][]string{
				"https://example.com": {"user1"},
			},
		},
		{
			name: "multiple visitors multiple URLs",
			operations: [][2]string{
				{"https://example.com", "user1"},
				{"https://example.com", "user2"},
				{"https://another.com", "user1"},
			},
			expectedData: map[string][]string{
				"https://example.com": {"user1", "user2"},
				"https://another.com": {"user1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &visitorRepoHandler{
				data: make(map[string]map[string]struct{}),
			}
			for _, op := range tt.operations {
				h.AddVisitor(op[0], op[1])
			}
			for url, expectedVisitors := range tt.expectedData {
				actualSet, exists := h.data[url]
				assert.True(t, exists, "url %q should exist", url)
				for _, visitor := range expectedVisitors {
					assert.Contains(t, actualSet, visitor, "visitor %q missing for url %q", visitor, url)
				}
			}
		})
	}
}

func TestVisitorHandler_CountVisitors(t *testing.T) {
	tests := []struct {
		name          string
		initialAdds   [][2]string
		queryURL      string
		expectedCount int
	}{
		{
			name:          "no visitors",
			initialAdds:   [][2]string{},
			queryURL:      "https://example.com",
			expectedCount: 0,
		},
		{
			name: "single visitor",
			initialAdds: [][2]string{
				{"https://example.com", "user1"},
			},
			queryURL:      "https://example.com",
			expectedCount: 1,
		},
		{
			name: "duplicate visitor",
			initialAdds: [][2]string{
				{"https://example.com", "user1"},
				{"https://example.com", "user1"},
			},
			queryURL:      "https://example.com",
			expectedCount: 1,
		},
		{
			name: "multiple visitors multiple URLs",
			initialAdds: [][2]string{
				{"https://example.com", "user1"},
				{"https://example.com", "user2"},
				{"https://another.com", "user1"},
			},
			queryURL:      "https://example.com",
			expectedCount: 2,
		},
		{
			name: "query URL with no visitors",
			initialAdds: [][2]string{
				{"https://example.com", "user1"},
			},
			queryURL:      "https://no-visitors.com",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &visitorRepoHandler{
				data: make(map[string]map[string]struct{}),
			}
			for _, op := range tt.initialAdds {
				h.AddVisitor(op[0], op[1])
			}
			got := h.CountVisitors(tt.queryURL)
			assert.Equal(t, tt.expectedCount, got, "CountVisitors(%q) should be %d", tt.queryURL, tt.expectedCount)
		})
	}
}

func TestVisitorHandler_ConcurrentAddCount(t *testing.T) {
	h := &visitorRepoHandler{
		data: make(map[string]map[string]struct{}),
	}

	var wg sync.WaitGroup
	url := "https://example.com"

	const (
		writeGoroutines = 50
		readGoroutines  = 20
	)
	for i := range writeGoroutines {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for j := range readGoroutines {
				visitorID := "user" + strconv.Itoa(g*readGoroutines+j)
				h.AddVisitor(url, visitorID)
			}
		}(i)
	}
	for range readGoroutines {
		wg.Go(func() {
			for range writeGoroutines {
				_ = h.CountVisitors(url)
			}
		})
	}
	wg.Wait()

	expected := writeGoroutines * readGoroutines
	got := h.CountVisitors(url)
	assert.Equal(t, expected, got, "final visitor count should match all added visitors")
}
