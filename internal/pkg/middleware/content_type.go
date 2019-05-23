package middleware

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

func ContentTypeMiddleware(inner routing.Handler) routing.Handler {
	return routing.Handler(func(ctx *routing.Context) error {
		ctx.SetContentType("application/json; charset=utf-8")
		return inner(ctx)
	})
}
