package middleware

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

type Middleware func(handler routing.Handler) routing.Handler

func ApplyMiddlewares(handler routing.Handler, middlewares ...Middleware) routing.Handler {
	ret := handler
	for _, mw := range middlewares {
		ret = mw(ret)
	}
	return ret
}
