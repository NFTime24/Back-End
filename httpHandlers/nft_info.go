package httpHandlers

import (
	"deukyunlee/protocol-camp/db"
	"encoding/json"
	"fmt"
	"net/http"
)

type NFTInfo struct {
	Id             int    `json:"id"`
	Artist_id      int    `json:"artist_id"`
	Price          int    `json:"price"`
	File_id        int    `json:"file_id"`
	Filesize       int    `json:"filesize"`
	Thumbnail_id   int    `json:"thumbnail_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	Owner_address  string `json:"owner_address"`
	Filename       string `json:"filename"`
	Filetype       string `json:"filetype"`
	Path           string `json:"path"`
	Thumbnail_path string `json:"thumbnail_path"`
	Artist_address string `json:"artist_address"`
	Artist_name    string `json:"artist_name"`
}

type NFTInfoBundle struct {
	NFTInfos []NFTInfo `json:"nft_infos"`
}

func GetNFTInfo(w http.ResponseWriter, r *http.Request) {
	nft_owner := r.URL.Query().Get("owner_address")

	db := db.DbConn()
	selDB, err := db.Query(fmt.Sprintf(`
	select wft.*,
	a.address as "artist_address", a.name as "artist_name"
	from
		(select wf.*,
		t.path as "thumnail_path"
		from
			(select 
			w.id, w.name, w.artist_id, 
			ifnull(w.price, 0) as "price", ifnull(w.description, "") as "description", 
			ifnull(w.category, "") as "category", ifnull(w.owner_address, "") as "owner_address", w.file_id,
			f.filename, f.filesize, f.filetype, f.path, ifnull(f.thumbnail_id, 0) as "thumbnail_id"
			from protocol_camp.work as w
			left join protocol_camp.file as f
			on w.file_id = f.id
			where w.owner_address = "%s") as wf
		left join protocol_camp.file as t
		on wf.thumbnail_id = t.id) as wft
	left join protocol_camp.artist as a
	on wft.artist_id = a.id
	`, nft_owner))
	if err != nil {
		panic(err.Error())
	}

	nftinfos := NFTInfoBundle{}
	for selDB.Next() {
		var id, artist_id, price, file_id, filesize, thumbnail_id int
		var name, description, category, owner_address, filename, filetype, path, thumbnail_path, artist_address, artist_name string
		//var fname, fsize, ftype, path string
		err = selDB.Scan(
			&id, &name, &artist_id, &price, &description, &category, &owner_address,
			&file_id, &filename, &filesize, &filetype, &path, &thumbnail_id, &thumbnail_path,
			&artist_address, &artist_name)
		if err != nil {
			panic(err.Error())
		}
		nftinfo := &NFTInfo{}

		nftinfo.Id = id
		nftinfo.Name = name
		nftinfo.Artist_id = artist_id
		nftinfo.Price = price
		nftinfo.Description = description
		nftinfo.Category = category
		nftinfo.Owner_address = owner_address
		nftinfo.File_id = file_id
		nftinfo.Filename = filename
		nftinfo.Filesize = filesize
		nftinfo.Filetype = filetype
		nftinfo.Path = path
		nftinfo.Thumbnail_id = thumbnail_id
		nftinfo.Thumbnail_path = thumbnail_path
		nftinfo.Artist_address = artist_address
		nftinfo.Artist_name = artist_name

		nftinfos.NFTInfos = append(nftinfos.NFTInfos, *nftinfo)
	}

	jData, err := json.Marshal(nftinfos)
	if err != nil {
		fmt.Printf("%s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jData)

	db.Close()
}
