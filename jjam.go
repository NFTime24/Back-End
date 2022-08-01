// package main

// import (
// 	"deukyunlee/protocol-camp/db"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/julienschmidt/httprouter"
// 	httpSwagger "github.com/swaggo/http-swagger"
// )

// type Artist struct {
// 	Aid      int
// 	Aname    string
// 	Aaddress string
// }

// type WorkInfo struct {
// 	Workid          int    `json: work_id`
// 	Workname        string `json: work_name`
// 	Workprice       int    `json: work_price`
// 	Workdescription string `json: work_description`
// 	Workcatogory    string `json: work_category`
// 	Artistname      string `json: artist_name`
// 	Artistaddress   string `json: artist_address`
// }

// func specificWork(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	workinfo := WorkInfo{}
// 	workinfos := []WorkInfo{}
// 	db := db.DbConn()
// 	w_name := r.URL.Query().Get("w_name")
// 	rows, err := db.Query(`select w.id, w.name, w.price, w.description, w.category, a.name, a.address  from work as w join artist as a where w.artist_id = a.id and w.name = ?;`, w_name)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for rows.Next() {
// 		var workid, workprice int
// 		var workname, workdescription, workcategory, artistname, artistaddress string

// 		err := rows.Scan(&workid, &workname, &workprice, &workdescription, &workcategory, &artistname, &artistaddress)
// 		if err != nil {
// 			panic(err)
// 		}
// 		workinfo.Workid = workid
// 		workinfo.Workname = workname
// 		workinfo.Workprice = workprice
// 		workinfo.Workdescription = workdescription
// 		workinfo.Workcatogory = workcategory
// 		workinfo.Artistname = artistname
// 		workinfo.Artistaddress = artistaddress

// 		workinfos = append(workinfos, workinfo)
// 	}
// 	jData, err := json.Marshal(workinfos)

// 	if err != nil {
// 		fmt.Printf("%s", err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(jData)
// 	db.Close()
// }

// func workInfoSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	workinfo := WorkInfo{}
// 	workinfos := []WorkInfo{}
// 	db := db.DbConn()
// 	rows, err := db.Query("select w.id, w.name, w.price, w.description, w.category, a.name, a.address from work as w join artist as a where w.artist_id = a.id")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for rows.Next() {
// 		var workid, workprice int
// 		var workname, workdescription, workcategory, artistname, artistaddress string

// 		err := rows.Scan(&workid, &workname, &workprice, &workdescription, &workcategory, &artistname, &artistaddress)
// 		if err != nil {
// 			panic(err)
// 		}
// 		workinfo.Workid = workid
// 		workinfo.Workname = workname
// 		workinfo.Workprice = workprice
// 		workinfo.Workdescription = workdescription
// 		workinfo.Workcatogory = workcategory
// 		workinfo.Artistname = artistname
// 		workinfo.Artistaddress = artistaddress

// 		workinfos = append(workinfos, workinfo)
// 	}
// 	jData, err := json.Marshal(workinfos)

// 	if err != nil {
// 		fmt.Printf("%s", err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(jData)
// 	db.Close()
// }
// func saveHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 	workname := r.FormValue("workname")
// 	artist := r.FormValue("artist")
// 	price := r.FormValue("price")
// 	description := r.FormValue("description")
// 	uploadFile, header, err := r.FormFile("upload_file")
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err)
// 		return
// 	}
// 	//fileByte = uploadFile.Open()

// 	defer uploadFile.Close()

// 	filename := header.Filename
// 	dirname := "./assets/uploadimage"
// 	os.MkdirAll(dirname, 0777)
// 	textname := strings.IndexByte(header.Filename, '.')
// 	extensionText := filename[textname:]
// 	if extensionText == ".mp4" {
// 		dirname = "./assets/uploadvideo"
// 	}
// 	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
// 	file, err := os.Create(filepath)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err)
// 		return
// 	}
// 	defer file.Close()
// 	fmt.Print(file.Name())

// 	io.Copy(file, uploadFile)
// 	w.WriteHeader(http.StatusOK)
// 	// fmt.Fprint(w, filepath)

// 	info, err := os.Stat(filepath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("filesize:", info.Size()/1024)
// 	filesize := int(info.Size() / 1024)
// 	//fmt.Println("type of filesize: ", reflect.TypeOf(filesize))
// 	fmt.Println("workname: ", workname)
// 	fmt.Println("artist: ", artist)
// 	fmt.Println("price: ", price)
// 	fmt.Println("description: ", description)

// 	text := strings.IndexByte(filename, '.')
// 	extension := filename[text:]
// 	fmt.Println("extension: ", extension)

// 	path := filepath[2:]
// 	fmt.Println("filepath: ", filepath[2:])

// 	fmt.Println("filename: ", filename)
// 	var filetype string
// 	switch extension {
// 	case ".png":
// 		filetype = "image/png"
// 	case ".jpg":
// 		filetype = "image/jpg"

// 	case ".jpeg":
// 		filetype = "image/jpeg"
// 	case ".gif":
// 		filetype = "image/gif"
// 	case ".mp4":
// 		filetype = "video/mp4"
// 	}

// 	db := db.DbConn()
// 	var id int
// 	var artistId int
// 	err = db.QueryRow("SELECT id FROM file order by id desc limit 1").Scan(&id)
// 	if err != nil {
// 		id = 0
// 		//log.Fatal(err)
// 	}
// 	//fmt.Println(id)
// 	id = id + 1
// 	insForm, err := db.Prepare("INSERT INTO file(id, filename, filesize, filetype, path) VALUES(?,?,?,?,?)")
// 	if err != nil {
// 		fmt.Println("file")
// 		panic(err.Error())
// 	}
// 	result, err := insForm.Exec(id, filename, filesize, filetype, path)
// 	if err != nil {
// 		fmt.Println("file")
// 		log.Fatal(err)
// 	} else {
// 		log.Println("data inserted successfully . . .")
// 	}
// 	n, err := result.RowsAffected()
// 	fmt.Println(n, "rows affected for file table")
// 	log.Printf("Successfully Uploaded File\n")

// 	err = db.QueryRow("SELECT id FROM work order by id desc limit 1").Scan(&artistId)
// 	if err != nil {
// 		id = 0
// 	}

// 	insFormWork, err := db.Prepare("INSERT INTO work(id, name, artist_id,price, description,category,file_id) VALUES(?,?,?,?,?,?,?)")
// 	if err != nil {
// 		fmt.Println("work")
// 		panic(err.Error())
// 	} else {
// 		log.Println("data insert successfully . . .")
// 	}

// 	resultWork, err := insFormWork.Exec(artistId+1, workname, artist, price, description, filetype, id)
// 	if err != nil {
// 		fmt.Println("work")
// 		log.Fatal(err)
// 	}
// 	num, err := resultWork.RowsAffected()
// 	fmt.Println(num, "rows affected for work table")
// 	log.Printf("Successfully Uploaded File\n")
// 	db.Close()
// 	fmt.Fprint(w, "file uploaded to"+filepath+"\n", num, " rows affected")
// }

// type NFTInfo struct {
// 	NftId             int    `json:"nft_id"`
// 	OwnerAddress      string `json:"owner_address"`
// 	WorkName          string `json:"work_name"`
// 	WorkPrice         int    `json:"work_price"`
// 	Description       string `json:"description"`
// 	Category          string `json:"category"`
// 	Filename          string `json:"filename"`
// 	Filesize          int    `json:"filesize"`
// 	Filetype          string `json:"filetype"`
// 	Path              string `json:"path"`
// 	ThumbnailPath     string `json:"thumbnail_path"`
// 	ArtistName        string `json:"artist_name"`
// 	ArtistProfilePath string `json:"artist_profile_path`
// 	ArtistAddress     string `json:"artist_address"`
// }

// // nft_id, work_price, filesize
// // owner_address, work_name, description, category, filename, filetype, path, thumbnail_path, artist_name, artist_profile_path, artist_address

// type NFTInfoBundle struct {
// 	NFTInfos []NFTInfo `json:"nft_infos"`
// }

// func getNFTInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	nft_owner := r.URL.Query().Get("owner_address")

// 	db := db.DbConn()
// 	selDB, err := db.Query(fmt.Sprintf(`
// 	select
// 		n.nft_id, n.owner_address, w.name as "work_name", w.price as "work_price", w.description, w.category,
// 		f.filename, f.filesize, f.filetype, f.path, t.path as "thumbnail_path",
// 		a.name as "artist_name", ifnull(p.path, "") as "artist_profile_path", a.address as "artist_address"
// 	from protocol_camp.nft as n
// 	left join protocol_camp.work as w on n.work_id = w.id
// 	left join protocol_camp.file as f on w.file_id = f.id
// 	left join protocol_camp.file as t on f.thumbnail_id = t.id
// 	left join protocol_camp.artist as a on w.artist_id = a.id
// 	left join protocol_camp.file as p on a.profile_id = p.id
// 	where n.owner_address = "%s"
// 	`, nft_owner))
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	nftinfos := NFTInfoBundle{}
// 	for selDB.Next() {
// 		var nft_id, work_price, filesize int
// 		var owner_address, work_name, description, category, filename, filetype, path, thumbnail_path, artist_name, artist_profile_path, artist_address string
// 		//var fname, fsize, ftype, path string
// 		err = selDB.Scan(&nft_id, &owner_address, &work_name, &work_price, &description, &category, &filename, &filesize, &filetype, &path,
// 			&thumbnail_path, &artist_name, &artist_profile_path, &artist_address)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		nftinfo := &NFTInfo{}

// 		nftinfo.NftId = nft_id
// 		nftinfo.OwnerAddress = owner_address
// 		nftinfo.WorkName = work_name
// 		nftinfo.WorkPrice = work_price
// 		nftinfo.Description = description
// 		nftinfo.Category = category
// 		nftinfo.Filename = filename
// 		nftinfo.Filesize = filesize
// 		nftinfo.Filetype = filetype
// 		nftinfo.Path = path
// 		nftinfo.ThumbnailPath = thumbnail_path
// 		nftinfo.ArtistName = artist_name
// 		nftinfo.ArtistProfilePath = artist_profile_path
// 		nftinfo.ArtistAddress = artist_address

// 		nftinfos.NFTInfos = append(nftinfos.NFTInfos, *nftinfo)
// 	}

// 	jData, err := json.Marshal(nftinfos)

// 	if err != nil {
// 		fmt.Printf("%s", err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(jData)
// 	db.Close()
// }

// func handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	fmt.Fprint(w, "Welcome to NFTime!\n")
// }

// func swaggerHandler(w http.ResponseWriter, r *http.Request) {
// 	swaggerFileUrl := "http://localhost:80/docs/swagger.json"
// 	handler := httpSwagger.Handler(httpSwagger.URL(swaggerFileUrl))
// 	handler.ServeHTTP(w, r)
// }

// // @title Swagger Example API
// // @version 1.0
// // @description This is a sample server Petstore server.
// // @termsOfService http://swagger.io/terms/

// // @contact.name API Support
// // @contact.url http://www.swagger.io/support
// // @contact.email support@swagger.io

// // @license.name Apache 2.0
// // @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// // @host petstore.swagger.io
// // @BasePath /v2

// func showImg(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
// 	// imgName := req.URL.Query().Get("name")
// 	// vars := mux.Vars(req)

// 	imgName := ps.ByName("imgName")
// 	fmt.Println(imgName)
// 	text := strings.IndexByte(imgName, '.')
// 	extension := imgName[text:]
// 	fmt.Println("extension: ", extension)
// 	var path string
// 	db := db.DbConn()
// 	selDB, err := db.Query("SELECT path FROM file where filename=?", imgName)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for selDB.Next() {
// 		err = selDB.Scan(&path)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	var filetype string
// 	switch extension {
// 	case ".png":
// 		filetype = "image/png"
// 	case ".jpg":
// 		filetype = "image/jpg"

// 	case ".jpeg":
// 		filetype = "image/jpeg"
// 	case ".gif":
// 		filetype = "image/gif"
// 	case ".mp4":
// 		filetype = "video/mp4"
// 	}
// 	if imgName != "" && path != "" {
// 		path = `./` + path

// 		fmt.Println(path)
// 		buf, err := ioutil.ReadFile(path)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		res.Header().Set("Content-Type", filetype)
// 		res.Write(buf)
// 	} else {
// 		panic(err)
// 	}
// }

// func main() {
// 	// r := chi.NewRouter()
// 	router := httprouter.New()

// 	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
// 	router.ServeFiles("/docs/*filepath", http.Dir("docs"))
// 	// http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
// 	// router.HandlerFunc(http.MethodGet, "/docs/index.html", swaggerHandler)
// 	// router.GET("/docs", swaggerHandler)
// 	// docs, err := chai.OpenAPI2(r)
// 	// if err != nil {
// 	// 	panic(fmt.Sprintf("failed to generate the swagger spec: %+v", err))
// 	// }
// 	// r.Get("/swagger/*", httpSwagger.Handler(
// 	// 	httpSwagger.URL("http://localhost:80/docs/swagger.json"), //The url pointing to API definition
// 	// ))
// 	// addCustomDocs(docs)
// 	// router.GET("/swagger", httpSwagger.WrapHandler)
// 	// router.Handler("/swagger", httpSwagger.WrapHandler)
// 	// // router.Handle("/swagger", httpSwagger.WrapHandler)
// 	// router.HandlerFunc("/swagger", httpSwagger.WrapHandler)

// 	//http.Handle("/swagger", httpSwagger.WrapHandler)
// 	router.GET("/assets/uploadimage/:imgName", showImg)
// 	// router.GET("/assets/uploadvideo/:videoName", showVideo)
// 	router.GET("/", handler)
// 	router.POST("/test", saveHandler)
// 	router.GET("/nft/specific", getNFTInfo)
// 	router.GET("/work/full", workInfoSearch)
// 	router.GET("/work/specific", specificWork)
// 	log.Fatal(http.ListenAndServe(":80", router))
// }