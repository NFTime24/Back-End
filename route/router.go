package route

import (
	"github.com/duke/api"

	"github.com/labstack/echo"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", api.Home)
	e.Static("/assets", "assets")
	e.GET("/users", api.GetFiles)
	e.GET("/works", api.GetWorks)
	e.GET("prepareAuth", api.PrepareAuth)
	e.GET("/getNFTInfo", api.GetNFTInfo)
	e.GET("/test", api.GetTest)
	e.GET("/file-upload", api.UploadWork)
	return e
}
