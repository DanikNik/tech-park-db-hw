package logger

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"log"
	"time"
)

// TODO: minimize logs for perfomance

func Logger(inner routing.Handler, name string) routing.Handler {
	return routing.Handler(func(ctx *routing.Context) error {
		start := time.Now()
		log.Printf("-----------[%s]-----------", ctx.RequestURI())

		err := inner(ctx)

		log.Printf(
			"%v %s %s %s",
			ctx.Response.StatusCode(),
			ctx.Method(),
			ctx.RequestURI(),
			time.Since(start),
		)
		return err
	})
}
