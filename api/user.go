package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

func PostUser(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["user_nickname"], params["user_address"], params["user_profile_id"])

	var id uint
	var user_id model.User
	user_nickname := params["user_nickname"]
	user_address := params["user_address"]
	user_profile_str := params["user_profile_id"]

	profile, _ := strconv.ParseUint(user_profile_str, 10, 32)
	user_profile_id := uint(profile)

	fmt.Println("user profile id:", user_profile_id)
	if user_profile_id == 0 {
		user_profile_id = 137
	}

	var checkAddress string
	db.Model(&user_id).Select("address").Where("address=?", user_address).Scan(&checkAddress)
	fmt.Println("checkAddress:", checkAddress)
	if len(checkAddress) == 0 {
		db.Model(&user_id).Select("id").Last(&id)
		id += 1
		fmt.Println("new id:", id)

		user_insert := model.User{ID: id, Address: user_address, Nickname: user_nickname, ProfileID: user_profile_id}

		db.Create(&user_insert)
	}

	type Result struct {
		Id        uint   `json:"id"`
		Address   string `json:"address"`
		Nickname  string `json:"nickname"`
		ProfileID uint   `json:"profile_id"`
		Path      string `json:"path"`
	}

	var result Result
	db.Select(`u.*, f.path`).
		Table("users as u").
		Joins("left join files as f on f.id = u.profile_id").
		Where("u.address=?", checkAddress).Scan(&result)

	fmt.Println(result)
	return c.JSON(http.StatusOK, result)
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
