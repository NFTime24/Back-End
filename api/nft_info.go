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
// @Param nft_id query string true "nft_id"
// @Router /getNFTInfoWithId [get]
func GetNFTInfoWithId(c echo.Context) error {
	// nft_owner := c.QueryParam("owner_address")
	nft_id_str := c.QueryParam("nft_id")
	nft_id, _ := strconv.ParseUint(nft_id_str, 10, 64)

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
	}

	db := db.ConnectDB()
	var users model.User
	var results Result
	db.Model(users).Select(`n.nft_id as nft_id,
	 w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category,f.filename as file_name, f.filesize as file_size, 
	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path, 
	 a.address as artist_address`).
		Joins("left join nfts as n on users.id = n.owner_id").
		Joins("left join works as w on n.works_id = w.work_id").
		Joins("left join files as f on w.file_id = f.id").
		Joins("left join files as t on f.thumbnail_id = t.id").
		Joins("left join artists as a on w.artist_id = a.id").
		Joins("left join files as p on a.profile_id = p.id").
		Where("n.nft_id=?", nft_id).Scan(&results)

	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}
