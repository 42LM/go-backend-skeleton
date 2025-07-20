package httpnone

import (
	"fmt"
	"net/http"
)

// HandlerFunc implements the std http.HandlerFunc for the none handler.
func (h *NoneHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Fprint(w, h.svc.FindNone(ctx))
}
