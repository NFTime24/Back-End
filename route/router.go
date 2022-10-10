package route

import (
	"github.com/duke/api"

	echo "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	e := echo.New()
	e.GET("/", api.Home)

	e.POST("/file-upload", api.UploadWork)
	e.Static("/assets", "assets")

	e.GET("/mintArtWithoutPaying", api.MintArtWithoutPaying)
	e.GET("/onSuccessKlip", api.OnSuccessKlip)
	e.GET("/mintToAddr", api.MintToAddr)
	e.GET("/addNFTWithWorkId", api.AddNFTWithWorkId)
	e.GET("/getNFTInfoWithId", api.GetNFTInfoWithId)

	e.POST("/user", api.PostUser)
	e.GET("/getUserWithAddress", api.GetUserWithAddress)

	e.POST("/workInfo", api.PostWork)
	e.GET("/getWorksInExhibition", api.GetWorksInExhibition)
	e.GET("/getWorksInfo", api.GetWorksInfoInExhibition)
	e.GET("/getWorkIdWithNftId", api.GetWorkIdWithNftId)
	e.GET("/getWorkInfoWithId", api.GetWorkInfoWithId)
	e.GET("/getTopWorks", api.GetTopWorks)
	e.GET("/getTopWorksWithCategory", api.GetTopWorksWithCategory)
	e.GET("/getTodayWorks", api.GetTodayWorks)
	e.GET("/getNewWorks", api.GetNewWorks)
	e.GET("/getAllWorks", api.GetAllWorks)
	e.GET("/getArtistWorks", api.GetArtistWorks)
	e.GET("/work/specific", api.GetSpecificWorkWithName)
	e.GET("/work/top10", api.GetTop10Works)

	e.POST("/artist", api.PostArtist)
	e.GET("/artist", api.ShowAllArtists)
	e.GET("/getActiveArtists", api.GetActiveArtists)
	e.GET("/getTopArtists", api.GetTopArtists)

	e.POST("/fantalk", api.PostFantalk)
	e.GET("/getArtistFantalks", api.GetArtistFantalks)

	e.POST("/like", api.UpdateLike)

	e.POST("/exhibition", api.PostExhibition)
	e.GET("/exhibition", api.ShowAllExhibitions)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/redirectTest", api.RedirectTest)
	return e
}
