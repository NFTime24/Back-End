package main

import (
	"fmt"
	"os"

	"github.com/duke/api"
	"github.com/duke/db"
	_ "github.com/duke/docs"
	"github.com/duke/route"
)

// @title NFTime Sample Swagger API
// @version 1.0
// @host 34.212.84.161
// @BasePath /
func main() {
	db.Init()
	e := route.Init()
	fmt.Println("Home: ", os.Getenv("APP_ENV"))

	api.KlipRequestMap = make(map[uint64]string)

	e.Logger.Fatal(e.Start(":80"))
}
