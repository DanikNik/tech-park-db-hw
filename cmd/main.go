package main

import (
	"fmt"
	"tech-park-db-hw/internal/pkg/db"
)

func main() {

	if err := db.Open(); err != nil {
		panic(err)
	}
	defer db.Close()
	//
	//router := customRouter.NewRouter()
	//log.Println("Server running at 5000")
	//panic(fasthttp.ListenAndServe(":5000", router.HandleRequest))

	user, _ := db.GetUser("qwerty")
	fmt.Println(*user)
}