package main

import (
	"github.com/duke/api"
	"github.com/duke/db"
	_ "github.com/duke/docs"
	"github.com/duke/route"
)

// @title NFTime Sample Swagger API
// @version 1.0
// @host localhost:80
// @BasePath /
func main() {
	db.Init()
	e := route.Init()
	api.KlipRKeyMap = make(map[uint64]string)

	e.Logger.Fatal(e.Start(":80"))
}
