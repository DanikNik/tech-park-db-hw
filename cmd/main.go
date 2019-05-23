package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"tech-park-db-hw/internal/pkg/db"
	customRouter "tech-park-db-hw/internal/pkg/router"
)

func main() {

	if err := db.Open(); err != nil {
		panic(err)
	}
	defer db.Close()

	router := customRouter.NewRouter()
	log.Println("Server running at 8080")
	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
