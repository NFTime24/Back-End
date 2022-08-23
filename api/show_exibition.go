package api

import (
	"fmt"
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary Get specific NFT
// @Description Get nft info
// @Tags exibition
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
		IsOwned              bool   `json:"is_owned"`
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
