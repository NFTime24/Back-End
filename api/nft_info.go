package api

import (
	"fmt"
	"net/http"
	"strconv"

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

type NFTInfo2 struct {
	NftId             int    `json:"nft_id"`
	WorkName          string `json:"work_name"`
	WorkPrice         int    `json:"work_price"`
	Description       string `json:"description"`
	Category          string `json:"category"`
	Filename          string `json:"filename"`
	Filesize          int    `json:"filesize"`
	Filetype          string `json:"filetype"`
	Path              string `json:"path"`
	ThumbnailPath     string `json:"thumbnail_path"`
	ArtistName        string `json:"artist_name"`
	ArtistProfilePath string `json:"artist_profile_path"`
	ArtistAddress     string `json:"artist_address"`
}

func GetNFTInfoWithId(c echo.Context) error {
	nft_id_str := c.QueryParam("nft_id")
	nft_id, _ := strconv.ParseUint(nft_id_str, 10, 64)
	db := db.DbConn()
	selDB, err := db.Query(fmt.Sprintf(`
	select 
		n.nft_id, w.name as "work_name", w.price as "work_price", w.description, w.category,
		f.filename, f.filesize, f.filetype, f.path, ifnull(t.path, "") as "thumbnail_path", 
		a.name as "artist_name", ifnull(p.path, "") as "artist_profile_path", a.address as "artist_address"
	from test.users as u
	left join test.nfts as n on u.id = n.owner_id
	left join test.works as w on n.works_id = w.work_id
	left join test.files as f on w.file_id = f.id
	left join test.files as t on f.thumbnail_id = t.id
	left join test.artists as a on w.artist_id = a.id
	left join test.files as p on a.profile_id = p.id
	where n.nft_id = %d
	`, nft_id))
	if err != nil {
		panic(err.Error())
	}

	nftinfo := &NFTInfo2{}
	for selDB.Next() {
		var nft_id, work_price, filesize int
		var work_name, description, category, filename, filetype, path, thumbnail_path, artist_name, artist_profile_path, artist_address string
		//var fname, fsize, ftype, path string
		err = selDB.Scan(&nft_id, &work_name, &work_price, &description, &category, &filename, &filesize, &filetype, &path,
			&thumbnail_path, &artist_name, &artist_profile_path, &artist_address)
		if err != nil {
			panic(err.Error())
		}
		nftinfo.NftId = nft_id
		nftinfo.WorkName = work_name
		nftinfo.WorkPrice = work_price
		nftinfo.Description = description
		nftinfo.Category = category
		nftinfo.Filename = filename
		nftinfo.Filesize = filesize
		nftinfo.Filetype = filetype
		nftinfo.Path = path
		nftinfo.ThumbnailPath = thumbnail_path
		nftinfo.ArtistName = artist_name
		nftinfo.ArtistProfilePath = artist_profile_path
		nftinfo.ArtistAddress = artist_address
	}

	// jData, err := json.Marshal(nftinfos)
	// if err != nil {
	// 	fmt.Printf("%s", err)
	// }

	// fmt.Printf("%s", string(jData))

	db.Close()
	return c.JSON(http.StatusOK, nftinfo)
}
