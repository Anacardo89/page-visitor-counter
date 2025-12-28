package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Anacardo89/page-visitor-counter/internal/repo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVisitorAPI_AddVisitor(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repo.NewMockVisitorRepo(ctrl)
	handler := NewVisitorAPI(mockRepo)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectCall     bool
	}{
		{
			name:           "valid request",
			requestBody:    `{"url":"https://example.com","visitor_id":"user1"}`,
			expectedStatus: http.StatusNoContent,
			expectCall:     true,
		},
		{
			name:           "missing url",
			requestBody:    `{"visitor_id":"user1"}`,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name:           "missing visitor_id",
			requestBody:    `{"url":"https://example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name:           "empty body",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name:           "duplicate visitor (ok)",
			requestBody:    `{"url":"https://example.com","visitor_id":"user1"}`,
			expectedStatus: http.StatusNoContent,
			expectCall:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockRepo.EXPECT().AddVisitor("https://example.com", "user1").Times(1)
			}
			req := httptest.NewRequest(http.MethodPost, "/visitors", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			handler.AddVisitor(w, req)
			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestVisitorAPI_CountVisitors(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repo.NewMockVisitorRepo(ctrl)
	handler := NewVisitorAPI(mockRepo)

	tests := []struct {
		name           string
		queryParam     string
		mockReturn     int
		expectedStatus int
		expectedBody   string
		expectCall     bool
	}{
		{
			name:           "valid request with visitors",
			queryParam:     "https://example.com",
			mockReturn:     5,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"visitor_count":5}` + "\n",
			expectCall:     true,
		},
		{
			name:           "valid request with zero visitors",
			queryParam:     "https://example.com",
			mockReturn:     0,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"visitor_count":0}` + "\n",
			expectCall:     true,
		},
		{
			name:           "missing url query param",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockRepo.EXPECT().CountVisitors(tt.queryParam).Return(tt.mockReturn).Times(1)
			}
			req := httptest.NewRequest(http.MethodGet, "/visitors?url="+tt.queryParam, nil)
			w := httptest.NewRecorder()
			handler.CountVisitors(w, req)
			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			if tt.expectCall {
				body, _ := io.ReadAll(res.Body)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}
