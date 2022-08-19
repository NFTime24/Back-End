package main

import (
	"fmt"
	"os"

	"github.com/duke/db"
	_ "github.com/duke/docs"
	"github.com/duke/route"
)

// @title NFTime Sample Swagger API
// @version 1.0
// @host localhost
// @BasePath /
func main() {
	db.Init()
	e := route.Init()
	fmt.Println("Home: ", os.Getenv("APP_ENV"))
	e.Logger.Fatal(e.Start(":80"))
}
