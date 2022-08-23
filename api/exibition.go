package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary Get specific NFT
// @Description Get nft info
// @Tags Exibition
// @Accept json
// @Produce json
// @Router /exibition [get]
func ShowAllExibitions(c echo.Context) error {

	type Result struct {
		ExibitonId           int    `json:"exibition_id"`
		ExibitionName        string `json:"exibition_name"`
		ExibitionDescription string `json:"exibition_description"`
		Link                 string `json:"link"`
		StartDate            string `json:"start_date"`
		EndDate              string `json:"end_date"`
		Filename             string `json:"filename"`
		FileSize             int    `json:"filesize"`
		FileType             string `json:"filetype"`
		FilePath             string `json:"path"`
	}

	db := db.DbManager()
	var exibitions model.Exibition
	var results []Result
	rows, err := db.Model(exibitions).Select(`exibitions.exibition_id as exibiton_id, exibitions.name as exibition_name, exibitions.description as exibition_description, 
    exibitions.start_date as start_date, exibitions.end_date as end_date, exibitions.link as link,f.filename,
    f.filename as file_name, f.filesize as file_size, 
	f.filetype as file_type, f.path as file_path`).
		Joins("left join files as f on exibitions.file_id = f.id").Rows()

	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
		fmt.Println(results[1])
	}

	return c.JSON(http.StatusOK, results)
}

// @Summary update exibition
// @Description update exibition
// @Tags Exibition
// @Accept json
// @Produce json
// @Param like body model.ExibitionCreateParam true "exibition data"
// @Router /exibition [post]
func PostExibition(c echo.Context) error {
	db := db.DbManager()

	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["name"], params["description"], params["start_date"], params["end_date"], params["file_id"], params["link"])

	name := params["name"]
	description := params["description"]
	start_date := params["start_date"]
	end_date := params["end_date"]
	file_str := params["file_id"]
	link := params["link"]
	file, _ := strconv.ParseUint(file_str, 10, 32)
	file_id := uint(file)
	var id uint
	var exibition_id model.Exibition
	db.Model(&exibition_id).Pluck("ExibitionID", &id)

	id += 1
	fmt.Println(id)
	exibition_insert := model.Exibition{ExibitionID: id, Name: name, Description: description, StartDate: start_date, EndDate: end_date, FileID: file_id, Link: link}
	db.Create(&exibition_insert)

	return c.JSON(http.StatusOK, params["name"])
}
