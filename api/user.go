package api

import (
	"fmt"
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

func PostUser(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["user_nickname"], params["user_address"])

	var id uint
	var user_id model.User
	user_nickname := params["user_nickname"]
	user_address := params["user_address"]

	db.Model(&user_id).Select("id").Last(&id)
	id += 1
	fmt.Println(id)

	user_insert := model.User{ID: id, Address: user_address, Nickname: user_nickname}

	db.Create(&user_insert)
	return c.JSON(http.StatusOK, params["user_address"])
}

func GetUserWithAddress(c echo.Context) error {
	address := c.QueryParam("address")

	type Result struct {
		Address  string `json:"address"`
		Nickname string `json:"nickname"`
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
