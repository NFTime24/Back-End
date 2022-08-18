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
// @Tags NFT
// @Accept json
// @Produce json
// @Param nft_id query string true "nft_id"
// @Router /getNFTInfoWithId [get]
func ShowAllExibitions(c echo.Context) error {

	type Result struct {
		ExibitonId          int    `json:"exibition_id"`
		ExibitionName       string `json:"exibition_name"`
		ExibitonDescription string `json:"exibition_description"`
		StartDate           string `json:"start_date"`
		EndDate             string `json:"end_date"`
		Filename            string `json:"filename"`
		FileSize            int    `json:"filesize"`
		FileType            string `json:"filetype"`
		FilePath            string `json:"path"`
	}

	db := db.DbManager()
	var exibitions model.Exibition
	var results Result
	db.Model(exibitions).Select(`exibitions.id as exibiton_id, exibitions.name as exibition_name, exibitions.description as exibition_description, 
    exibitions.start_date as start_date, exibitions.end_date as end_date, f.filename,
    f.filename as file_name, f.filesize as file_size, 
	f.filetype as file_type, f.path as file_path`).
		Joins("left join files as f on exibitions.file_id = f.id").Scan(&results)

	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}