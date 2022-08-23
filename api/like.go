package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary update like
// @Description update like
// @Tags Like
// @Accept json
// @Produce json
// @Param like body model.LikeCreateParam true "like data"
// @Router /like [post]
func UpdateLike(c echo.Context) (err error) {
	db := db.DbManager()
	type Like struct {
		UserId uint `json:"user_id"`
		WorkId uint `json:"work_id"`
	}
	like := Like{}

	err = c.Bind(&like)
	if err != nil {
		log.Printf("Failed processing update like request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	fmt.Println(like.UserId)
	var likes model.Like
	like_index := -1
	type This []struct {
		UserId uint `json:"user_id"`
		WorkId uint `json:"work_id"`
	}

	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
	db.Model(likes).Select("like_id").
		Where("owner_id=? and works_id=?", like.UserId, like.WorkId).Scan(&like_index)

	if like_index == -1 {
		// 좋아요가 되어있지 않은 상태
		var id uint
		var likeIs model.Like

		db.Model(&likeIs).Pluck("LikeID", &id)
		// fmt.Println(id)
		id += 1
		updateLike := model.Like{LikeID: id, OwnerID: like.UserId, WorksID: like.WorkId}
		err := db.Create(&updateLike)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "err occured")
		}
		return c.String(http.StatusOK, "Successfully updated")
	} else {
		// 좋아요를 한 상태라면 삭제
		db.Where("like_id =?", like_index).Delete(&like)
		return c.String(http.StatusOK, "Successfully deleted")
	}
}
