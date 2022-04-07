package resServer

import (
	"UResNode/internal/logic/resServer"
	"net/http"

	"UResNode/internal/svc"
	"UResNode/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UploadResHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadResReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := resServer.NewUploadResLogic(r.Context(), ctx)
		resp, err := l.UploadRes(req, r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
