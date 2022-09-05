package api

import (
	"fmt"
	"net/http"

	"github.com/duke/db"
	"github.com/labstack/echo/v4"
)

func GetUserWithAddress(c echo.Context) error {
	address := c.QueryParam("address")

	type Result struct {
		Address  string `json:"address"`
		NickName string `json:"nickname"`
		Path     string `json:"path"`
	}

	db := db.DbManager()
	var result Result
	db.Select(`u.address, u.nickname, f.path`).
		Table("users as u").
		Joins("left join files as f on f.id = u.profile_id").
		Where("u.address=?", address).Scan(&result)

	fmt.Println(result)
	return c.JSON(http.StatusOK, result)
}
