package api

import (
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

func GetFiles(c echo.Context) error {
	db := db.DbManager()

	files := []model.File{}
	db.Find(&files)

	return c.JSON(http.StatusOK, files)
}
