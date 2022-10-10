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
// @Tags exhibition
// @Accept json
// @Produce json
// @Router /exhibition [get]
func ShowAllExhibitions(c echo.Context) error {

	type Result struct {
		ExibitonId            int    `json:"exhibition_id"`
		ExhibitionName        string `json:"exhibition_name"`
		ExhibitionDescription string `json:"exhibition_description"`
		Link                  string `json:"link"`
		StartDate             string `json:"start_date"`
		EndDate               string `json:"end_date"`
		Filename              string `json:"filename"`
		FileSize              int    `json:"filesize"`
		FileType              string `json:"filetype"`
		FilePath              string `json:"path"`
	}

	db := db.DbManager()
	var exhibitions model.Exhibition
	var results []Result
	rows, err := db.Model(exhibitions).Select(`exhibitions.exhibition_id as exibiton_id, exhibitions.name as exhibition_name, exhibitions.description as exhibition_description, 
    exhibitions.start_date as start_date, exhibitions.end_date as end_date, exhibitions.link as link,f.filename,
    f.filename as file_name, f.filesize as file_size, 
	f.filetype as file_type, f.path as file_path`).
		Joins("left join files as f on exhibitions.file_id = f.id").Rows()

	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
		fmt.Println(results[1])
	}

	return c.JSON(http.StatusOK, results)
}

// @Summary update exhibition
// @Description update exhibition
// @Tags exhibition
// @Accept json
// @Produce json
// @Param like body model.ExhibitionCreateParam true "exhibition data"
// @Router /exhibition [post]
func PostExhibition(c echo.Context) error {
	db := db.DbManager()

	params := make(map[string]string)
	bind_params := c.Bind(&params)
	fmt.Println(bind_params)

	name := params["name"]
	description := params["description"]
	start_date := params["start_date"]
	end_date := params["end_date"]
	file_str := params["file_id"]
	link := params["link"]
	file, _ := strconv.ParseUint(file_str, 10, 32)
	file_id := uint(file)
	var exibition_id uint
	var exhibition_dbId model.Exhibition
	db.Model(&exhibition_dbId).Pluck("ExhibitionID", &exibition_id)

	exibition_id += 1
	fmt.Println(exibition_id)
	exhibition_insert := model.Exhibition{ExhibitionID: exibition_id, Name: name, Description: description, StartDate: start_date, EndDate: end_date, FileID: file_id, Link: link}
	db.Create(&exhibition_insert)

	return c.JSON(http.StatusOK, params["name"])
}
