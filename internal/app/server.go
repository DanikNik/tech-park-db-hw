package app

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"tech-park-db-hw/internal/pkg/router"
)

type Service struct {
	Port   string
	Router *routing.Router
}

func StartApp() error {
	service := Service{
		Port:   ":8080",
		Router: router.NewRouter(),
	}
	log.Printf("Server running at %v\n", service.Port)
	return fasthttp.ListenAndServe(service.Port, service.Router.HandleRequest)
}
