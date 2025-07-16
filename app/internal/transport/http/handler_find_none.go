package http

import (
	"fmt"
	"net/http"

	"go-backend-skeleton/app/internal/svc"
)

// TODO: How to handle logging?
func handlerFindNone(svc svc.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Fprint(w, svc.FindNone(ctx))
	})
}
