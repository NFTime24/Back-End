package route

import (
	"github.com/duke/api"

	echo "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	e := echo.New()
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.GET("/", api.Home)
	e.Static("/assets", "assets")
	e.GET("/users", api.GetFiles)
	e.GET("/works", api.GetWorks)
	e.GET("prepareAuth", api.PrepareAuth)
	e.GET("/mintArt", api.MintArt)
	e.GET("/getKlipResult", api.GetKlipResult)
	e.GET("/getNFTInfo", api.GetNFTInfo)
	e.GET("/getNFTInfoWithId", api.GetNFTInfoWithId)
	e.GET("/test", api.GetTest)
	// e.GET("/test2", api.TestGo)
	e.GET("/file-upload", api.UploadWork)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
