package route

import (
	"github.com/duke/api"

	echo "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	e := echo.New()
	e.GET("/", api.Home)
	e.Static("/assets", "assets")
	e.GET("/work/specific", api.GetSpecificWorkWithName)
	e.GET("/work/top10", api.GetTop10Works)
	e.GET("/mintArt", api.MintArt)
	e.GET("/getNFTInfoWithId", api.GetNFTInfoWithId)
	e.GET("/getWorkIdWithNftId", api.GetWorkIdWithNftId)
	e.GET("/getWorkInfoWithId", api.GetWorkInfoWithId)
	e.POST("/like", api.UpdateLike)
	e.GET("/exibition", api.ShowAllExibitions)
	e.POST("/exibition", api.PostExibition)
	e.POST("/artist", api.PostArtist)
	e.POST("/workInfo", api.PostWork)
	e.POST("/file-upload", api.UploadWork)
	e.GET("/getWorksInfo", api.GetWorksInfoInExibition)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
