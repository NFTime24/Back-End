package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary Get specific NFT
// @Description Get nft info
// @Tags NFT
// @Accept json
// @Produce json
// @Param ex_id query string true "ex_id"
// @Param user_id query string true "user_id"
// @Router /getWorksInfo [get]
func GetWorksInfoInExibition(c echo.Context) error {

	// nft_owner := c.QueryParam("owner_address")
	ex_id_str := c.QueryParam("ex_id")
	ex_id, _ := strconv.ParseUint(ex_id_str, 10, 64)

	user_id_str := c.QueryParam("user_id")
	user_id, _ := strconv.ParseUint(user_id_str, 10, 64)

	type Result struct {
		NftId         int    `json:"nft_id"`
		WorkName      string `json:"work_name"`
		Price         int    `json:"work_price"`
		Description   string `json:"description"`
		WorkCategory  string `json:"category"`
		FileName      string `json:"filename"`
		FileSize      int    `json:"filesize"`
		FileType      string `json:"filetype"`
		FilePath      string `json:"path"`
		ThumbnailPath string `json:"thumbnail_path"`
		ArtistName    string `json:"artist_name"`
		ProfilePath   string `json:"artist_profile_path"`
		ArtistAddress string `json:"artist_address"`
		UserId        uint   `json:"user_id"`
		ExibitionId   uint   `json:"exibition_id"`
		IsOwned       bool   `json:"is_owned"`
	}

	db := db.DbManager()
	var users model.User
	var results []Result
	var results_user []Result

	rows_user, err := db.Model(users).Select(`n.nft_id as nft_id, w.exibitions_id as exibition_id,
    w.name as work_name, w.price as price, w.description as description,
    w.category as work_category,f.filename as file_name, f.filesize as file_size,
    f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
    a.address as artist_address, users.id as user_id`).
		Joins("left join nfts as n on users.id = n.owner_id").
		Joins("left join works as w on n.works_id = w.work_id").
		Joins("left join files as f on w.file_id = f.id").
		Joins("left join files as t on f.thumbnail_id = t.id").
		Joins("left join artists as a on w.artist_id = a.id").
		Joins("left join files as p on a.profile_id = p.id").
		Joins("left join exibitions as e on e.exibition_id = w.exibitions_id").
		Where("w.exibitions_id=? and users.id=?", ex_id, user_id).Rows()
	if err != nil {
		panic(err)
	}
	for rows_user.Next() {
		db.ScanRows(rows_user, &results_user)
	}

	fmt.Println(results_user)
	rows, err := db.Model(users).Select(`n.nft_id as nft_id, w.exibitions_id as exibition_id,
	 w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category,f.filename as file_name, f.filesize as file_size, 
	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path, 
	 a.address as artist_address, users.id as user_id`).
		Joins("left join nfts as n on users.id = n.owner_id").
		Joins("left join works as w on n.works_id = w.work_id").
		Joins("left join files as f on w.file_id = f.id").
		Joins("left join files as t on f.thumbnail_id = t.id").
		Joins("left join artists as a on w.artist_id = a.id").
		Joins("left join files as p on a.profile_id = p.id").
		Joins("left join exibitions as e on e.exibition_id = w.exibitions_id").
		Where("w.exibitions_id=?", ex_id).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
	}
	// grayscale := "_grayscale"
	for key, _ := range results {
		for key, _ := range results_user {
			if results[key].WorkName == results_user[key].WorkName {
				results[key].IsOwned = true
			}
			fmt.Println(key)
		}
		fmt.Println(results[key].IsOwned)
	}

	return c.JSON(http.StatusOK, results)
}

func GetWorkInfoWithId(c echo.Context) error {

	// nft_owner := c.QueryParam("owner_address")
	work_id_str := c.QueryParam("work_id")
	work_id, _ := strconv.ParseUint(work_id_str, 10, 64)

	type Result struct {
		WorkId        int    `json:"work_id"`
		WorkName      string `json:"work_name"`
		Price         int    `json:"work_price"`
		Description   string `json:"description"`
		WorkCategory  string `json:"category"`
		FileName      string `json:"filename"`
		FileSize      int    `json:"filesize"`
		FileType      string `json:"filetype"`
		FilePath      string `json:"path"`
		ThumbnailPath string `json:"thumbnail_path"`
		ArtistName    string `json:"artist_name"`
		ProfilePath   string `json:"artist_profile_path"`
		ArtistAddress string `json:"artist_address"`
	}

	db := db.DbManager()
	var results Result
	db.Select(`w.work_id as work_id,
	 w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category, f.filename as file_name, f.filesize as file_size, 
	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path, 
	 a.address as artist_address`).
		Table("works as w").
		Joins("left join files as f on w.file_id = f.id").
		Joins("left join files as t on f.thumbnail_id = t.id").
		Joins("left join artists as a on w.artist_id = a.id").
		Joins("left join files as p on a.profile_id = p.id").
		Where("work_id=?", work_id).Scan(&results)

	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}
