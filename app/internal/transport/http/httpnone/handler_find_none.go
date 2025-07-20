package httpnone

import (
	"fmt"
	"net/http"
)

// TODO: Middleware for logging
func (h *NoneHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Fprint(w, h.svc.FindNone(ctx))
}
