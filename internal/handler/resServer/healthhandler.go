package resServer

import (
	"UResNode/internal/logic/resServer"
	"net/http"

	"UResNode/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func HealthHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := resServer.NewHealthLogic(r.Context(), ctx)
		err := l.Health()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
