package httpmsg

import (
	"fmt"
	"net/http"
)

// HandlerFunc implements the std http.HandlerFunc for the message handler.
func (h *MsgHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	fmt.Fprint(w, h.svc.FindMsg(ctx, id))
}
