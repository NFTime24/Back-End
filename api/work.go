package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo"
)

func GetWorks(c echo.Context) error {
	db := db.ConnectDB()

	name := c.QueryParam("name")
	var artists model.Artist
	// db.Where
	db.Joins("JOIN works w on w.work_id = artists.id").
		Where("w.name =?", name).Find(&artists)
	return c.JSON(http.StatusOK, artists)
}

func UploadWork(c echo.Context) error {
	db := db.ConnectDB()

	// name := c.FormValue("name")
	// workname := c.FormValue("workname")
	// artist := c.FormValue("artist")
	// price := c.FormValue("price")
	// description := c.FormValue("description")
	file, err := c.FormFile("upload_file")
	if err != nil {
		panic(err)
	}
	filename := file.Filename
	dirname := "assets/uploadimage"
	os.MkdirAll(dirname, 0777)
	src, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	textname := strings.IndexByte(file.Filename, '.')
	extensionText := filename[textname:]
	if extensionText == ".mp4" {
		dirname = "assets/uploadvideo"
		os.MkdirAll(dirname, 0777)
	}
	filepath := fmt.Sprintf("%s/%s", dirname, filename)
	fileTemp := fmt.Sprintf("%s", dirname)
	text := strings.IndexByte(filename, '.')
	extension := filename[text:]
	var filetype string
	var tempFile *os.File
	// dst, err := os.Create(filepath)
	switch extension {
	case ".png":
		tempFile, err = ioutil.TempFile(fileTemp, "upload-*.png")
		filetype = "image/png"
	case ".jpg":
		tempFile, err = ioutil.TempFile(fileTemp, "upload-*.jpg")
		filetype = "image/jpg"
	case ".jpeg":
		tempFile, err = ioutil.TempFile(fileTemp, "upload-*.jpeg")
		filetype = "image/jpeg"
	case ".mp4":
		tempFile, err = ioutil.TempFile(fileTemp, "upload-*.mp4")
		filetype = "video/mp4"
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("this:", tempFile.Name())

	defer tempFile.Close()
	fmt.Println(filepath)

	if _, err = io.Copy(tempFile, src); err != nil {
		return err
	}

	fmt.Println(filetype)
	tempPath := tempFile.Name()
	path := filepath[2:]
	filesize := uint(file.Size)
	fmt.Println("filepath: ", tempPath)
	fmt.Println("path: ", path)
	fmt.Println("filename: ", filename)
	fmt.Println("filetype: ", filetype)
	fmt.Println("filesize: ", filesize)

	var id uint
	var files model.File

	db.Model(&files).Pluck("Id", &id)

	id += 1

	insertFile := model.File{ID: id, Filename: filename, Filesize: filesize, Filetype: filetype, Path: tempPath}
	db.Create(&insertFile)

	return c.JSON(http.StatusOK, insertFile)
}
