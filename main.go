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

type Page struct {
	Title string
	Body  []byte
}

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
	// fmt.Println(Artist{Aname: "sircoon"})
	// fmt.Println(artists)
	// fmt.Println(artist)
	// fmt.Println(artist.Aid[0])
	fmt.Println(artists[0].Aid)
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
	// fmt.Println(Artist{Aname: "sircoon"})
	// fmt.Println(artists)
	// fmt.Println(artist)
	// fmt.Println(artist.Aid[0])
	fmt.Println(artists[0].Aid)
	defer db.Close()

	// filename := title + ".txt"
	// body, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	return nil, err
	// }
	return &Artist{Aid: artists[0].Aid, Aname: artists[0].Aname, Aaddress: artists[0].Aaddress}, nil
}

// 사용자가 해당 페이지를 볼 수 있도록 함.

// func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
// 	p, err := loadPage(title)
// 	if err != nil {
// 		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
// 		return
// 	}
// 	renderTemplate(w, "view", p)
// }

//

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
		//fmt.Println("image/png")
	case ".jpg":
		filetype = "image/jpg"
		//fmt.Println("image/jpg")
	case ".jpeg":
		filetype = "image/jpeg"
		//fmt.Println("image/jpeg")
	case ".gif":
		filetype = "image/gif"
		//fmt.Println("image/gif")
	case ".mp4":
		filetype = "video/mp4"
		//fmt.Println("video/mp4")
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

	insForm, err := db.Prepare("INSERT INTO file(id, filename, filesize, filetype, path) VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println("file")
		panic(err.Error())
	} else {
		log.Println("data insert successfully . . .")
	}
	result, err := insForm.Exec(id+1, filename, filesize, filetype, path)
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
		//log.Fatal(err)
	}
	// err = db.QueryRow("SELECT id FROM artist where name = ?",).Scan(&id)
	// if err != nil {
	// 	id = 0
	// 	// //log.Fatal(err)
	// }

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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

func getNFTInfo(w http.ResponseWriter, r *http.Request) {
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
