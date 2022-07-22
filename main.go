package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func dbConn() (db *sql.DB) {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er.Error())
	}
	dbDriver := os.Getenv("DB_Driver")
	dbUser := os.Getenv("DB_User")
	dbPass := os.Getenv("DB_Password")
	dbName := os.Getenv("DB_Name")
	dbHost := os.Getenv("DB_HOST")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type upfile struct {
	ID    int
	Fname string
	Fsize string
	Ftype string
	Path  string
	Count int
}

// var tmplIndex = *template.Template
var tmpl = template.Must(template.ParseGlob("templates/uploadfile.html"))

// var tmplIndex = template.Must(template.ParseGlob("templates/index.html"))

func upload(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT id,filename,filesize,filetype,path FROM file ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	upld := upfile{}
	res := []upfile{}
	for selDB.Next() {
		var id int
		var fname, fsize, ftype, path string
		err = selDB.Scan(&id, &fname, &fsize, &ftype, &path)
		if err != nil {
			panic(err.Error())
		}
		upld.ID = id
		upld.Fname = fname
		upld.Fsize = fsize
		upld.Ftype = ftype
		upld.Path = path
		res = append(res, upld)

	}

	upld.Count = len(res)

	if upld.Count > 0 {
		tmpl.ExecuteTemplate(w, "uploadfile.html", res)
	} else {
		tmpl.ExecuteTemplate(w, "uploadfile.html", nil)
	}
	db.Close()
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	var id int

	r.ParseMultipartForm(200000)
	if r == nil {
		fmt.Fprintf(w, "No files can be selected\n")
	}

	formdata := r.MultipartForm

	fil := formdata.File["files"]

	selDB, err := db.Query("SELECT id FROM file ORDER BY id DESC limit 1")
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&id)
		if err != nil {
			id = 0
		}
	}
	id = id + 1
	for i := range fil {

		file, err := fil[i].Open()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer file.Close()

		fname := fil[i].Filename
		fsize := fil[i].Size
		kilobytes := fsize / 1024
		// megabytes := (float64)(kilobytes / 1024) // cast to type float64

		ftype := fil[i].Header.Get("Content-type")
		var tempFile *os.File

		text := strings.IndexByte(fname, '.')
		extension := fname[text:]

		switch extension {
		case ".png":
			tempFile, err = ioutil.TempFile("assets/uploadimage", "upload-*.png")
		case ".jpg":
			tempFile, err = ioutil.TempFile("assets/uploadimage", "upload-*.jpg")
		case ".jpeg":
			tempFile, err = ioutil.TempFile("assets/uploadimage", "upload-*.jpeg")
		case ".mp4":
			tempFile, err = ioutil.TempFile("assets/uploadvideo", "upload-*.mp4")
		}

		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		filepath := tempFile.Name()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		insForm, err := db.Prepare("INSERT INTO file(id,filename, filesize, filetype, path) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		} else {
			log.Println("data insert successfully . . .")
		}
		insForm.Exec(id, fname, kilobytes, ftype, filepath)

		log.Printf("Successfully Uploaded File\n")
		defer db.Close()

		http.Redirect(w, r, "/", 301)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")

	// selDB, err := db.Query("SELECT path from upload where id=?", emp)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println(selDB)
	// err3 := os.Remove(selDB)
	fmt.Println(emp)
	delForm, err := db.Prepare("DELETE FROM file WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("deleted successfully")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er.Error())
	}
	port := os.Getenv("PORT")
	log.Println("Server started on PORT: " + port)

	http.HandleFunc("/delete", delete)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", upload)
	http.HandleFunc("/uploadfiles", uploadFiles)
	http.Handle("/test", http.FileServer(http.Dir("/assets")))
	// http.HandleFunc("/showimg", showImg)
	// http.Handle("/", http.FileServer(http.Dir("assets/uploadimage")))
	http.ListenAndServe(":80", nil)
}

// func index(w http.ResponseWriter, r *http.Request) {
// 	// io.WriteString(w, "Hello fcc")
// 	tmplIndex.ExecuteTemplate(w, "index.html", nil)
// }

//sudo nohup go run main.go &

// package main

// import (
// 	"io"
// 	"net/http"
// )

// func main() {
// 	http.HandleFunc("/", index)
// 	http.ListenAndServe(":80", nil)
// }

// func index(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "Hello Sircoon!")
// }

// func showImg(res http.ResponseWriter, req *http.Request) {
// 	imgName := req.URL.Query().Get("name")
// 	var path string
// 	db := dbConn()
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

// 	if imgName != "" && path != "" {
// 		path = `./` + path

// 		fmt.Println(path)
// 		buf, err := ioutil.ReadFile(path)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		res.Header().Set("Content-Type", "image/png")
// 		res.Write(buf)
// 	} else {
// 		panic(err)
// 	}

// }
