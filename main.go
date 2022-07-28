package main

import (
	"deukyunlee/protocol-camp/db"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Artist struct {
	Aid      int
	Aname    string
	Aaddress string
}

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// Page의 Body 부분을 text file로 저장

// func (p *Page) save() error {
// 	filename := p.Title + ".txt"
// 	return ioutil.WriteFile(filename, p.Body, 0600)
// }

func dbSelect() []Artist {
	artist := Artist{}
	artists := []Artist{}
	db := db.DbConn()
	rows, err := db.Query("select id, name, address from artist")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var name, address string
		err := rows.Scan(&id, &name, &address)
		if err != nil {
			panic(err)
		}
		artist.Aid = id
		artist.Aname = name
		artist.Aaddress = address
		artists = append(artists, artist)
	}
	defer db.Close()
	return artists
}

// title 변수를 통해 파일이름을 생성한 후 파일의 내용을 읽어들여 Page literal에 대한 ptr 반환

func loadPage() (*Artist, error) {
	artist := Artist{}
	artists := []Artist{}
	db := db.DbConn()
	rows, err := db.Query("select id, name, address from artist")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var name, address string
		err := rows.Scan(&id, &name, &address)
		if err != nil {
			panic(err)
		}
		artist.Aid = id
		artist.Aname = name
		artist.Aaddress = address
		artists = append(artists, artist)
	}

	fmt.Println(artists[0].Aid)
	defer db.Close()

	return &Artist{Aid: artists[0].Aid, Aname: artists[0].Aname, Aaddress: artists[0].Aaddress}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	table := dbSelect()

	err := templates.ExecuteTemplate(w, "edit.html", table)
	if err != nil {
		panic(err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	workname := r.FormValue("workname")
	fmt.Println(workname)
	artist := r.FormValue("artist")
	price := r.FormValue("price")
	description := r.FormValue("description")
	uploadFile, header, err := r.FormFile("upload_file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	//fileByte = uploadFile.Open()

	defer uploadFile.Close()

	filename := header.Filename
	dirname := "./assets/uploadimage"
	os.MkdirAll(dirname, 0777)
	textname := strings.IndexByte(header.Filename, '.')
	extensionText := filename[textname:]
	if extensionText == ".mp4" {
		dirname = "./assets/uploadvideo"
	}
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	defer file.Close()
	fmt.Print(file.Name())

	io.Copy(file, uploadFile)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath)

	info, err := os.Stat(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("filesize:", info.Size()/1024)
	filesize := int(info.Size() / 1024)
	//fmt.Println("type of filesize: ", reflect.TypeOf(filesize))
	fmt.Println("workname: ", workname)
	fmt.Println("artist: ", artist)
	fmt.Println("price: ", price)
	fmt.Println("description: ", description)

	text := strings.IndexByte(filename, '.')
	extension := filename[text:]
	fmt.Println("extension: ", extension)

	path := filepath[2:]
	fmt.Println("filepath: ", filepath[2:])

	fmt.Println("filename: ", filename)
	var filetype string
	switch extension {
	case ".png":
		filetype = "image/png"
	case ".jpg":
		filetype = "image/jpg"

	case ".jpeg":
		filetype = "image/jpeg"
	case ".gif":
		filetype = "image/gif"
	case ".mp4":
		filetype = "video/mp4"
	}

	db := db.DbConn()
	var id int
	var artistId int
	err = db.QueryRow("SELECT id FROM file order by id desc limit 1").Scan(&id)
	if err != nil {
		id = 0
		//log.Fatal(err)
	}
	//fmt.Println(id)
	id = id + 1
	insForm, err := db.Prepare("INSERT INTO file(id, filename, filesize, filetype, path) VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println("file")
		panic(err.Error())
	} else {
		log.Println("data insert successfully . . .")
	}
	result, err := insForm.Exec(id, filename, filesize, filetype, path)
	if err != nil {
		fmt.Println("file")
		log.Fatal(err)
	}
	n, err := result.RowsAffected()
	fmt.Println(n, "rows affected")
	log.Printf("Successfully Uploaded File\n")

	err = db.QueryRow("SELECT id FROM work order by id desc limit 1").Scan(&artistId)
	if err != nil {
		id = 0
	}

	insFormWork, err := db.Prepare("INSERT INTO work(id, name, artist_id,price, description,category,file_id) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("work")
		panic(err.Error())
	} else {
		log.Println("data insert successfully . . .")
	}

	resultWork, err := insFormWork.Exec(artistId+1, workname, artist, price, description, filetype, id)
	if err != nil {
		fmt.Println("work")
		log.Fatal(err)
	}
	num, err := resultWork.RowsAffected()
	fmt.Println(num, "rows affected")
	log.Printf("Successfully Uploaded File\n")
	db.Close()
	// p := &Page{Title: title, Body: []byte(body)}
	// err := p.save()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseGlob("./templates/edit.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

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
	ArtistProfilePath string `json:"artist_profile_path`
	ArtistAddress     string `json:"artist_address"`
}

// nft_id, work_price, filesize
// owner_address, work_name, description, category, filename, filetype, path, thumbnail_path, artist_name, artist_profile_path, artist_address

type NFTInfoBundle struct {
	NFTInfos []NFTInfo `json:"nft_infos"`
}

func getNFTInfo(w http.ResponseWriter, r *http.Request) {
	nft_owner := r.URL.Query().Get("owner_address")

	db := db.DbConn()
	selDB, err := db.Query(fmt.Sprintf(`
	select 
		n.nft_id, n.owner_address, w.name as "work_name", w.price as "work_price", w.description, w.category,
		f.filename, f.filesize, f.filetype, f.path, t.path as "thumbnail_path", 
		a.name as "artist_name", ifnull(p.path, "") as "artist_profile_path", a.address as "artist_address"
	from protocol_camp.nft as n
	left join protocol_camp.work as w on n.work_id = w.id
	left join protocol_camp.file as f on w.file_id = f.id
	left join protocol_camp.file as t on f.thumbnail_id = t.id
	left join protocol_camp.artist as a on w.artist_id = a.id
	left join protocol_camp.file as p on a.profile_id = p.id
	where n.owner_address = "%s"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jData)
	db.Close()
}

func main() {
	dbSelect()
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	//http.HandleFunc("/view", upload)

	//http.Handle("/edit/*", http.StripPrefix("/edit/view", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/getNFTInfo", getNFTInfo)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	//http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("./css/")))

	log.Fatal(http.ListenAndServe(":80", nil))
}
