package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary update work
// @Description update work
// @Tags Work
// @Accept json
// @Produce json
// @Param like body model.WorkInfoCreateParam true "work info data"
// @Router /workInfo [post]
func PostWork(c echo.Context) error {
	db := db.DbManager()
	params := make(map[string]string)
	test := c.Bind(&params)
	fmt.Println(test)
	fmt.Println(params["work_name"], params["work_price"], params["work_description"], params["work_category"], params["file_id"], params["artist_id"])

	var id uint
	var work_id model.Work
	work_name := params["work_name"]
	work_price_str := params["work_price"]
	work_description := params["work_description"]
	work_category := params["work_category"]

	file_id_str := params["file_id"]
	artist_id_str := params["artist_id"]

	price, _ := strconv.ParseUint(work_price_str, 10, 32)
	work_price := uint(price)

	file, _ := strconv.ParseUint(file_id_str, 10, 32)
	file_id := uint(file)

	artist, _ := strconv.ParseUint(artist_id_str, 10, 32)
	artist_id := uint(artist)
	fmt.Println(artist_id)
	db.Model(&work_id).Pluck("work_id", &id)
	id += 1
	fmt.Println(id)

	work_insert := model.Work{WorkID: id, Name: work_name, Price: work_price, Description: work_description, Category: work_category, FileID: file_id, ArtistID: artist_id}

	db.Create(&work_insert)
	return c.JSON(http.StatusOK, params["work_name"])
}

// @Summary Get specific NFT
// @Description Get nft info
// @Tags NFT
// @Accept json
// @Produce json
// @Param ex_id query string true "ex_id"
// @Param user_id query string true "user_id"
// @Router /getWorksInfo [get]
func GetWorksInfoInExhibition(c echo.Context) error {

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
		ExhibitionId  uint   `json:"exhibition_id"`
		IsOwned       bool   `json:"is_owned"`
	}

	db := db.DbManager()
	var users model.User
	var results []Result
	var results_user []Result

	rows_user, err := db.Model(users).Select(`n.nft_id as nft_id, w.exhibitions_id as exhibition_id,
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
		Joins("left join exhibitions as e on e.exhibition_id = w.exhibitions_id").
		Where("w.exhibitions_id=? and users.id=?", ex_id, user_id).Rows()
	if err != nil {
		panic(err)
	}
	for rows_user.Next() {
		db.ScanRows(rows_user, &results_user)
	}

	fmt.Println(results_user)
	rows, err := db.Model(users).Select(`n.nft_id as nft_id, w.exhibitions_id as exhibition_id,
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
		Joins("left join exhibitions as e on e.exhibition_id = w.exhibitions_id").
		Where("w.exhibitions_id=?", ex_id).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
	}

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

// @Summary get specific work
// @Description Get works
// @Tags Work
// @Accept json
// @Produce json
// @Param name query string true "name"
// @Router /work/specific [get]
func GetSpecificWorkWithName(c echo.Context) error {
	name := c.QueryParam("name")
	// 구조체 멤버변수 이름과 DB에서 가져오는 컬럼명이 일치해야함
	type Result struct {
		WorkName        string
		ArtistName      string
		WorkDescription string
	}
	db := db.DbManager()

	// var artists model.Artist
	var works model.Work
	var results Result

	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
	db.Model(works).Select("works.name as work_name, works.description as work_description, artists.name as artist_name").Joins("left join artists on works.work_id = artists.id").Where("works.name=?", name).Scan(&results)
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}

// @Summary get top 10 works
// @Description get top 10 works
// @Tags Work
// @Accept json
// @Produce json
// @Router /work/top10 [get]
func GetTop10Works(c echo.Context) error {
	// 구조체 멤버변수 이름과 DB에서 가져오는 컬럼명이 일치해야함
	// filepath, workname, artistname
	type Result struct {
		WorkName   string
		ArtistName string
		FilePath   string
	}
	db := db.DbManager()

	// var artists model.Artist
	var works model.Work
	var results []Result

	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
	rows, err := db.Model(works).Select("works.name as work_name, f.path as file_path, a.name as artist_name").
		Joins("left join files as f on works.file_id = f.id").
		Joins("left join artists as a on works.artist_id = a.id").Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &results)
	}
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}

// @Summary Get specific Work
// @Description Get work info in Exibition
// @Tags NFT
// @Accept json
// @Produce json
// @Param ex_id query string true "ex_id"
// @Router /getWorksInExhibition [get]
func GetWorksInExhibition(c echo.Context) error {

	// nft_owner := c.QueryParam("owner_address")
	ex_id_str := c.QueryParam("ex_id")
	ex_id, _ := strconv.ParseUint(ex_id_str, 10, 64)

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
		ExhibitionId  uint   `json:"exhibition_id"`
	}

	db := db.DbManager()
	var users model.User
	var results []Result

	rows, err := db.Model(users).Select(`w.work_id as work_id, w.exhibitions_id as exhibition_id,
	 w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category,f.filename as file_name, f.filesize as file_size, 
	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path, 
	 a.address as artist_address, users.id as user_id`).
		Joins("left join files as f on w.file_id = f.id").
		Joins("left join files as t on f.thumbnail_id = t.id").
		Joins("left join artists as a on w.artist_id = a.id").
		Joins("left join files as p on a.profile_id = p.id").
		Joins("left join exhibitions as e on e.exhibition_id = w.exhibitions_id").
		Where("w.exhibitions_id=?", ex_id).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		db.ScanRows(rows, &results)
	}

	return c.JSON(http.StatusOK, results)
}
