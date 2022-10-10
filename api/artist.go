package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary update artist
// @Description update artist
// @Tags Artist
// @Accept json
// @Produce json
// @Param like body model.ArtistCreateParam true "artist data"
// @Router /artist [post]
func PostArtist(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)

	var artist_id uint
	var artist_dbId model.Artist
	artist_name := params["artist_name"]
	artist_address := params["artist_address"]
	artist_profile_str := params["artist_profile_id"]

	profile, _ := strconv.ParseUint(artist_profile_str, 10, 32)
	artist_profile_id := uint(profile)

	db.Model(&artist_dbId).Pluck("ID", &artist_id)
	artist_id += 1
	fmt.Println(artist_id)

	artist_insert := model.Artist{ID: artist_id, Name: artist_name, Address: artist_address, ProfileID: artist_profile_id}

	db.Create(&artist_insert)
	return c.JSON(http.StatusOK, params)
}

// @Summary artist info
// @Description Get All Artist Info
// @Tags Artist
// @Accept json
// @Produce json
// @Router /artist [get]
func ShowAllArtists(c echo.Context) error {

	type Result struct {
		Id         uint   `json:"id"`
		Name       string `json:"name"`
		Address    string `json:"address"`
		Profile_id uint   `json:"profile_id"`
	}

	db := db.DbManager()
	var artists model.Artist
	var results []Result
	rows, err := db.Model(artists).Select(`artists.*`).Rows()

	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
	}

	return c.JSON(http.StatusOK, results)
}
