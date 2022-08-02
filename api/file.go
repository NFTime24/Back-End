package api

import (
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo"
)

func GetFiles(c echo.Context) error {
	db := db.ConnectDB()

	files := []model.File{}
	db.Find(&files)

	return c.JSON(http.StatusOK, files)
}
