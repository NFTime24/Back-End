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
	e.GET("/work/specific", api.GetSpecificWork)
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
	// e.GET("/getNFTInfo", api.GetNFTInfo)
	// e.GET("/test", api.GetTest)
	// e.GET("/test2", api.TestGo)
	e.POST("/file-upload", api.UploadWork)
	e.GET("/getWorksInfo", api.GetWorksInfoInExibition)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
