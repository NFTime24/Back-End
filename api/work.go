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
	c.Request().ParseMultipartForm(256 >> 20)

	thumb_file, err := c.FormFile("thumb_file")
	if err != nil {
		panic(err)
	}
	thumb_filename := thumb_file.Filename
	thumb_dirname := "assets/uploadimage"
	dirname := "assets/uploadimage"
	os.MkdirAll(thumb_dirname, 0777)

	thumb_src, err := thumb_file.Open()
	if err != nil {
		panic(err)
	}
	defer thumb_src.Close()

	thumb_filepath := fmt.Sprintf("%s/%s", thumb_dirname, thumb_filename)
	thumb_fileTemp := fmt.Sprintf("%s", thumb_dirname)
	thumb_text := strings.IndexByte(thumb_filename, '.')
	thumb_extension := thumb_filename[thumb_text:]
	var thumb_filetype string
	var thumb_tempFile *os.File
	// dst, err := os.Create(filepath)
	switch thumb_extension {
	case ".png":
		thumb_tempFile, err = ioutil.TempFile(thumb_fileTemp, "upload-*.png")
		thumb_filetype = "image/png"
	case ".jpg":
		thumb_tempFile, err = ioutil.TempFile(thumb_fileTemp, "upload-*.jpg")
		thumb_filetype = "image/jpg"
	case ".jpeg":
		thumb_tempFile, err = ioutil.TempFile(thumb_fileTemp, "upload-*.jpeg")
		thumb_filetype = "image/jpeg"
	case ".mp4":
		panic("img thumbnail only")
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("this:", thumb_tempFile.Name())

	defer thumb_tempFile.Close()
	fmt.Println(thumb_filepath)

	if _, err = io.Copy(thumb_tempFile, thumb_src); err != nil {
		return err
	}

	thumb_tempPath := thumb_tempFile.Name()
	thumb_filesize := uint(thumb_file.Size)

	var thumb_id uint
	var thumb_files model.File

	db.Model(&thumb_files).Pluck("Id", &thumb_id)

	thumb_id += 1

	thumb_insertFile := model.File{ID: thumb_id, Filename: thumb_filename, Filesize: thumb_filesize, Filetype: thumb_filetype, Path: thumb_tempPath}
	db.Create(&thumb_insertFile)

	file, err := c.FormFile("upload_file")
	if err != nil {
		panic(err)
	}
	filename := file.Filename

	src, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	textname := strings.IndexByte(filename, '.')
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
	var thumbnail_id uint
	var files model.File

	db.Model(&files).Pluck("Id", &id)
	db.Model(&files).Pluck("Id", &thumbnail_id)
	id += 1

	insertFile := model.File{ID: id, Filename: filename, Filesize: filesize, Filetype: filetype, Path: tempPath, ThumbnailID: &thumbnail_id}
	db.Create(&insertFile)

	return c.JSON(http.StatusOK, insertFile)
}
