package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

func PostArtist(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["artist_name"], params["artist_address"], params["artist_profile_id"])

	var id uint
	var artist_id model.Artist
	artist_name := params["artist_name"]
	artist_address := params["artist_address"]
	artist_profile_str := params["artist_profile_id"]

	profile, _ := strconv.ParseUint(artist_profile_str, 10, 32)
	artist_profile_id := uint(profile)

	db.Model(&artist_id).Pluck("ID", &id)
	id += 1
	fmt.Println(id)

	artist_insert := model.Artist{ID: id, Name: artist_name, Address: artist_address, ProfileID: artist_profile_id}

	db.Create(&artist_insert)
	return c.JSON(http.StatusOK, params["artist_name"])
}
