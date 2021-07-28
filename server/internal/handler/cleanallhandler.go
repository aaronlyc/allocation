package handler

import (
	"net/http"

	"github.com/aaronlyc/allocation/server/internal/logic"
	"github.com/aaronlyc/allocation/server/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func CleanAllHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCleanAllLogic(r.Context(), ctx)
		err := l.CleanAll()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
