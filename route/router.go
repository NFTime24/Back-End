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
	e.GET("/mintArtWithoutPaying", api.MintArtWithoutPaying)
	e.GET("/onSuccessKlip", api.OnSuccessKlip)
	e.GET("/mintToAddr", api.MintToAddr)
	e.GET("/addNFTWithWorkId", api.AddNFTWithWorkId)
	e.GET("/getNFTInfoWithId", api.GetNFTInfoWithId)
	e.GET("/getWorkIdWithNftId", api.GetWorkIdWithNftId)
	e.GET("/getWorkInfoWithId", api.GetWorkInfoWithId)
	e.GET("/getTopWorks", api.GetTopWorks)
	e.GET("/getTopWorksWithCategory", api.GetTopWorksWithCategory)
	e.GET("/getTodayWorks", api.GetTodayWorks)
	e.GET("/getNewWorks", api.GetNewWorks)
	e.GET("/getAllWorks", api.GetAllWorks)
	e.POST("/like", api.UpdateLike)
	e.GET("/exhibition", api.ShowAllExhibitions)
	e.POST("/exhibition", api.PostExhibition)
	e.POST("/user", api.PostUser)
	e.POST("/artist", api.PostArtist)
	e.POST("/workInfo", api.PostWork)
	e.POST("/file-upload", api.UploadWork)
	e.GET("/getWorksInExhibition", api.GetWorksInExhibition)
	e.GET("/getWorksInfo", api.GetWorksInfoInExhibition)
	e.GET("/artist", api.ShowAllArtists)
	e.GET("/getActiveArtists", api.GetActiveArtists)
	e.GET("/getUserWithAddress", api.GetUserWithAddress)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/redirectTest", api.RedirectTest)
	return e
}
