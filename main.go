package main

import (
	"github.com/duke/route"
)

func main() {
	//db.Init()
	e := route.Init()

	e.Logger.Fatal(e.Start(":80"))
}
