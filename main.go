package main

import (
	"github.com/duke/api"
	"github.com/duke/db"
	"github.com/duke/route"
)

func main() {
	db.Init()
	e := route.Init()
	api.KlipRKeyMap = make(map[uint64]string)

	e.Logger.Fatal(e.Start(":80"))
}
