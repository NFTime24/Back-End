package main

import (
	"deukyunlee/protocol-camp/db"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"deukyunlee/protocol-camp/httpHandlers"

	_ "github.com/go-sql-driver/mysql"
)

type Page struct {
	Title string
	Body  []byte
}

// Page의 Body 부분을 text file로 저장

// func (p *Page) save() error {
// 	filename := p.Title + ".txt"
// 	return ioutil.WriteFile(filename, p.Body, 0600)
// }

// title 변수를 통해 파일이름을 생성한 후 파일의 내용을 읽어들여 Page literal에 대한 ptr 반환

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
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
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	workname := r.FormValue("workname")
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
	err = db.QueryRow("SELECT id FROM file order by id desc limit 1").Scan(&id)
	if err != nil {
		id = 0
		//log.Fatal(err)
	}
	//fmt.Println(id)

	insForm, err := db.Prepare("INSERT INTO file(id, filename, filesize, filetype, path) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	} else {
		log.Println("data insert successfully . . .")
	}
	result, err := insForm.Exec(id+1, filename, filesize, filetype, path)
	if err != nil {
		log.Fatal(err)
	}
	n, err := result.RowsAffected()
	fmt.Println(n, "rows affected")
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

func main() {

	_, err := os.Stat(filepath.Join(".", "assets/stylesheets", "test.css"))
	if err != nil {
		panic(err)
	}

	// file2, err := os.ReadFile("./uploads/짱구.jpeg")
	// fmt.Println(string(file2))
	// fmt.Println(file2)

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	//http.HandleFunc("/view", upload)

	//http.Handle("/edit/*", http.StripPrefix("/edit/view", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/getNFTInfo", httpHandlers.GetNFTInfo)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	//http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("./css/")))

	log.Fatal(http.ListenAndServe(":80", nil))
}
