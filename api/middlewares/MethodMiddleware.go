package middlewares

import (
	"io"
	"net/http"
	"strings"
)

// MethodMiddleware blocks non specified methods
func MethodMiddleware(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for i := range methods {
			if strings.EqualFold(methods[i], r.Method) {
				next(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Invalid Method")
	}
}
