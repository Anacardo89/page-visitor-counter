package repo

type VisitorRepo interface {
	AddVisitor(url, visitorID string)
	CountVisitors(url string) int
}
