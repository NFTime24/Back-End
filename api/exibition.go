package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

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
