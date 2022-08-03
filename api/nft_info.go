package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duke/db"
	"github.com/labstack/echo"
)

type NFTInfo struct {
	NftId             int    `json:"nft_id"`
	OwnerAddress      string `json:"owner_address"`
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

// nft_id, work_price, filesize
// owner_address, work_name, description, category, filename, filetype, path, thumbnail_path, artist_name, artist_profile_path, artist_address

type NFTInfoBundle struct {
	NFTInfos []NFTInfo `json:"nft_infos"`
}

func GetNFTInfo(c echo.Context) error {
	nft_owner := c.QueryParam("owner_address")
	db := db.DbConn()
	selDB, err := db.Query(fmt.Sprintf(`
	select 
		n.nft_id, u.address as "owner_address", w.name as "work_name", w.price as "work_price", w.description, w.category,
		f.filename, f.filesize, f.filetype, f.path, t.path as "thumbnail_path", 
		a.name as "artist_name", ifnull(p.path, "") as "artist_profile_path", a.address as "artist_address"
	from test.users as u
	left join test.nfts as n on u.id = n.owner_id
	left join test.works as w on n.works_id = w.work_id
	left join test.files as f on w.file_id = f.id
	left join test.files as t on f.thumbnail_id = t.id
	left join test.artists as a on w.artist_id = a.id
	left join test.files as p on a.profile_id = p.id
	where u.id = "%s"
	`, nft_owner))
	if err != nil {
		panic(err.Error())
	}

	nftinfos := NFTInfoBundle{}
	for selDB.Next() {
		var nft_id, work_price, filesize int
		var owner_address, work_name, description, category, filename, filetype, path, thumbnail_path, artist_name, artist_profile_path, artist_address string
		//var fname, fsize, ftype, path string
		err = selDB.Scan(&nft_id, &owner_address, &work_name, &work_price, &description, &category, &filename, &filesize, &filetype, &path,
			&thumbnail_path, &artist_name, &artist_profile_path, &artist_address)
		if err != nil {
			panic(err.Error())
		}
		nftinfo := &NFTInfo{}

		nftinfo.NftId = nft_id
		nftinfo.OwnerAddress = owner_address
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

		nftinfos.NFTInfos = append(nftinfos.NFTInfos, *nftinfo)
	}

	jData, err := json.Marshal(nftinfos)
	if err != nil {
		fmt.Printf("%s", err)
	}

	db.Close()
	return c.JSON(http.StatusOK, jData)
}
