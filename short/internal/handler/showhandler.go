package handler

import (
	"net/http"

	"short/internal/logic"
	"short/internal/svc"
	"short/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 参数解析
		var req types.ShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 参数校验 validate
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Errorw(
				"validator check failed",
				logx.LogField{Key: "err", Value: err.Error()},
			)

			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShowLogic(r.Context(), svcCtx)
		resp, err := l.Show(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 返回重定向
			http.Redirect(w, r, resp.LongUrl, http.StatusFound)
		}
	}
}
