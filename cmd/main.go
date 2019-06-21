package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/router"
)

func main() {

	if err := db.Open(); err != nil {
		panic(err)
	}
	defer db.Close()
	//
	serverRouter := router.NewRouter()
	log.Println("Server running at 5000")
	panic(fasthttp.ListenAndServe(":5000", serverRouter.HandleRequest))

	//user, _ := db.GetUser("qwerty")
	//fmt.Println(*user)
}
