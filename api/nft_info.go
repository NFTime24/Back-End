package api

import (
	"fmt"
	"net/http"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary NFT
// @Description Get nft info
// @Accept json
// @Produce json
// @Param owner_address query string true "NFT owner_address"
// @Router /getNFTInfo [get]
func GetNFTInfo(c echo.Context) error {
	nft_owner := c.QueryParam("owner_address")

	type Result struct {
		NftId         int
		OwnerAddress  string
		WorkName      string
		Price         int
		Description   string
		WorkCategory  string
		FileName      string
		FileSize      int
		FileType      string
		FilePath      string
		ThumbnailPath string
		ArtistName    string
		ProfilePath   string
		ArtistAddress string
	}

	db := db.ConnectDB()
	var users model.User
	var results Result
	db.Model(users).Select(`n.nft_id as nft_id, users.address as owner_address,
	 w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category,f.filename as file_name, f.filesize as file_size, 
	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_paht, a.name as artist_name, p.path as profile_path, 
	 a.address as artist_address`).
		Joins("left join nfts as n on users.id = n.owner_id").
		Joins("left join works as w on n.works_id = w.work_id").
		Joins("left join files as f on w.file_id = f.id").
		Joins("left join files as t on f.thumbnail_id = t.id").
		Joins("left join artists as a on w.artist_id = a.id").
		Joins("left join files as p on a.profile_id = p.id").Where("users.address=?", nft_owner).Scan(&results)

	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}
