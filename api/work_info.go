package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

func PostWork(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["work_name"], params["work_price"], params["work_description"], params["work_category"], params["file_id"], params["artist_id"])

	var id uint
	var work_id model.Work
	work_name := params["work_name"]
	work_price_str := params["work_price"]
	work_description := params["work_description"]
	work_category := params["work_category"]

	file_id_str := params["file_id"]
	artist_id_str := params["artist_id"]

	price, _ := strconv.ParseUint(work_price_str, 10, 32)
	work_price := uint(price)

	file, _ := strconv.ParseUint(file_id_str, 10, 32)
	file_id := uint(file)

	artist, _ := strconv.ParseUint(artist_id_str, 10, 32)
	artist_id := uint(artist)
	fmt.Println(artist_id)
	db.Model(&work_id).Pluck("work_id", &id)
	id += 1
	fmt.Println(id)

	work_insert := model.Work{WorkID: id, Name: work_name, Price: work_price, Description: work_description, Category: work_category, FileID: file_id, ArtistID: artist_id}

	db.Create(&work_insert)
	return c.JSON(http.StatusOK, params["work_name"])
}
