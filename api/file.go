package api

import (
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary upload file
// @Description upload file and thumbnail
// @Tags File
// @Accept json
// @Produce json
// @Param upload_file formData file true "file you want to upload"
// @Router /file-upload [post]
func GetFiles(c echo.Context) error {
	db := db.DbManager()

	files := []model.File{}
	db.Find(&files)

	return c.JSON(http.StatusOK, files)
}
