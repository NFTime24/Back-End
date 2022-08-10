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
	e.GET("/mintArt", api.MintArt)
	e.GET("/getKlipResult", api.GetKlipResult)
	e.GET("/getNFTInfo", api.GetNFTInfo)
	e.GET("/getNFTInfoWithId", api.GetNFTInfoWithId)
	e.GET("/test", api.GetTest)
	// e.GET("/test2", api.TestGo)
	e.GET("/file-upload", api.UploadWork)
	return e
}
