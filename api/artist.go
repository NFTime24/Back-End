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
		fmt.Println(results[1])
	}

	return c.JSON(http.StatusOK, results)
}

func GetActiveArtists(c echo.Context) error {
	type Result struct {
		Id         uint   `json:"id"`
		Name       string `json:"name"`
		Address    string `json:"address"`
		Profile_id uint   `json:"profile_id"`
	}

	db := db.DbManager()
	var results []Result
	rows, err := db.Select(`a.*`).
		Table(`artists as a`).
		Joins("join works as w on w.artist_id = a.id").
		Group("a.id").Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
	}
	return c.JSON(http.StatusOK, results)
}

func GetTopArtists(c echo.Context) error {
	type Result struct {
		Id        uint   `json:"id"`
		Name      string `json:"name"`
		Address   string `json:"address"`
		ProfileId uint   `json:"profile_id"`
		Path      string `json:"path"`
	}

	db := db.DbManager()
	var results []Result
	rows, err := db.Select(`a.*, f.path`).
		Table(`(
			Select w.artist_id, count(w.artist_id) as count
			from test.nfts as n
			left join test.works as w on w.work_id = n.works_id
			group by w.artist_id
			order by count desc ) as base`).
		Joins("left join artists as a on base.artist_id = a.id").
		Joins("left join files as f on f.id = a.profile_id").Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
	}
	return c.JSON(http.StatusOK, results)
}
