package handler

import (
	"net/http"

	"UResNode/internal/logic/test"
	"UResNode/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func TestRefreshTokenHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewTestRefreshTokenLogic(r.Context(), ctx)
		err := l.TestRefreshToken()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
