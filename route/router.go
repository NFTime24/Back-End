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
	e.GET("/addNFTWithWorkId", api.AddNFTWithWorkId)
	e.GET("/getNFTInfoWithId", api.GetNFTInfoWithId)
	e.GET("/getWorkIdWithNftId", api.GetWorkIdWithNftId)
	e.GET("/getWorkInfoWithId", api.GetWorkInfoWithId)
	e.GET("/getTopWorks", api.GetTopWorks)
	e.POST("/like", api.UpdateLike)
	e.GET("/exhibition", api.ShowAllExhibitions)
	e.POST("/exhibition", api.PostExhibition)
	e.POST("/artist", api.PostArtist)
	e.POST("/workInfo", api.PostWork)
	e.POST("/file-upload", api.UploadWork)
	e.GET("/getWorksInExhibition", api.GetWorksInExhibition)
	e.GET("/getWorksInfo", api.GetWorksInfoInExhibition)
	e.GET("/artist", api.ShowAllArtists)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e
}
