package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

func PostFantalk(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["artist_id"], params["owner_id"], params["post_text"])

	var id uint
	var fantalk_id model.Fantalk
	artist_id, _ := strconv.ParseUint(params["artist_id"], 10, 32)
	owner_id, _ := strconv.ParseUint(params["owner_id"], 10, 32)
	post_text := params["post_text"]

	db.Model(&fantalk_id).Select("post_id").Last(&id)
	id += 1
	fmt.Println("new id:", id)

	creative_time := time.Now()
	modify_time := time.Now()

	fantalk_insert := model.Fantalk{
		Post_id:    id,
		ArtistID:   uint(artist_id),
		OwnerID:    uint(owner_id),
		PostText:   post_text,
		LikeCount:  0,
		CreateTime: &creative_time,
		ModifyTime: &modify_time,
	}

	db.Create(&fantalk_insert)

	return c.String(http.StatusOK, strconv.FormatUint(uint64(id), 10))
}
