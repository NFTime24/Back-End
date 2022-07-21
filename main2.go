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
