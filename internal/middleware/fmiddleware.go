package middleware

import "net/http"

type FMiddleware struct {
}

func NewFMiddleware() *FMiddleware {
	return &FMiddleware{}
}

func (m *FMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
