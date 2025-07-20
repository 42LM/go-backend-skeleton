package httpmsg

import (
	"fmt"
	"net/http"
)

func (h *MsgHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	fmt.Fprint(w, h.svc.FindMsg(ctx, id))
}
