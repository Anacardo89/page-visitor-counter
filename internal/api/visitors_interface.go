package api

import "net/http"

type VisitorAPI interface {
	AddVisitor(w http.ResponseWriter, r *http.Request)
	CountVisitors(w http.ResponseWriter, r *http.Request)
}
