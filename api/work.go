package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

// @Summary get specific work
// @Description Get works
// @Tags Work
// @Accept json
// @Produce json
// @Param name query string true "name"
// @Router /work/specific [get]
func GetSpecificWork(c echo.Context) error {
	name := c.QueryParam("name")
	// 구조체 멤버변수 이름과 DB에서 가져오는 컬럼명이 일치해야함
	type Result struct {
		WorkName        string
		ArtistName      string
		WorkDescription string
	}
	db := db.DbManager()

	// var artists model.Artist
	var works model.Work
	var results Result

	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
	db.Model(works).Select("works.name as work_name, works.description as work_description, artists.name as artist_name").Joins("left join artists on works.work_id = artists.id").Where("works.name=?", name).Scan(&results)
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}

// @Summary get top 10 works
// @Description get top 10 works
// @Tags Work
// @Accept json
// @Produce json
// @Router /work/top10 [get]
func GetTop10Works(c echo.Context) error {
	// 구조체 멤버변수 이름과 DB에서 가져오는 컬럼명이 일치해야함
	// filepath, workname, artistname
	type Result struct {
		WorkName   string
		ArtistName string
		FilePath   string
	}
	db := db.ConnectDB()

	// var artists model.Artist
	var works model.Work
	var results []Result

	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
	rows, err := db.Model(works).Select("works.name as work_name, f.path as file_path, a.name as artist_name").
		Joins("left join files as f on works.file_id = f.id").
		Joins("left join artists as a on works.artist_id = a.id").Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &results)
	}
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}

func UploadWork(c echo.Context) error {
	db := db.ConnectDB()
	err := c.Request().ParseMultipartForm(512 >> 20)
	if err != nil {
		panic(err)
	}

	thumb_dirname := "assets/uploadimage"
	dirname := "assets/uploadimage"
	var thumb_filetype string
	var thumb_tempFile *os.File
	var filetype string
	var tempFile *os.File
	if err != nil {
		log.Println(err.Error())
		return err
	}
	thumb_file, err := c.FormFile("thumbnail_file")
	files := c.Request().MultipartForm.File["thumbnail_file"]
	if err != nil {
		panic(err)
	}
	thumb_filename := thumb_file.Filename
	os.MkdirAll(thumb_dirname, 0777)
	thumb_fileTemp := fmt.Sprintf("%s", thumb_dirname)
	thumb_text := strings.IndexByte(thumb_filename, '.')
	thumb_extension := thumb_filename[thumb_text:]

	fmt.Println(thumb_extension)
	for _, file := range files {
		thumb_file, err := file.Open()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		fmt.Println(thumb_file)
		defer thumb_file.Close()
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
			thumb_tempFile, err = ioutil.TempFile(thumb_fileTemp, "upload-*.mp4")
			thumb_filetype = "video/mp4"
		}
		if err != nil {
			panic(err)
		}

		if _, err = io.Copy(thumb_tempFile, thumb_file); err != nil {
			return err
		}

		defer thumb_tempFile.Close()

		defer thumb_file.Close()

	}
	thumb_tempPath := thumb_tempFile.Name()
	thumb_filesize := uint(thumb_file.Size)

	var thumb_id uint
	var thumb_files model.File

	db.Model(&thumb_files).Pluck("Id", &thumb_id)

	thumb_id += 1

	thumb_insertFile := model.File{ID: thumb_id, Filename: thumb_filename, Filesize: thumb_filesize, Filetype: thumb_filetype, Path: thumb_tempPath}
	db.Create(&thumb_insertFile)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	file, err := c.FormFile("upload_file")
	upload_files := c.Request().MultipartForm.File["upload_file"]
	if err != nil {
		panic(err)
	}

	filename := file.Filename

	text := strings.IndexByte(filename, '.')
	extension := filename[text:]
	fmt.Println(extension)
	if extension == ".mp4" {
		dirname = "assets/uploadvideo"
		fmt.Println("here")
		os.MkdirAll(dirname, 0777)
	}
	fileTemp := fmt.Sprintf("%s", dirname)
	fmt.Println(extension)

	for _, file2 := range upload_files {
		file, err := file2.Open()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		fmt.Println(file2)
		defer file.Close()
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

		if _, err = io.Copy(tempFile, file); err != nil {
			return err
		}

		defer tempFile.Close()

		defer file.Close()

	}
	tempPath := tempFile.Name()
	filesize := uint(file.Size)

	var id uint
	var filesIs model.File

	db.Model(&filesIs).Pluck("Id", &id)

	id += 1

	insertFile := model.File{ID: id, Filename: filename, Filesize: filesize, Filetype: filetype, Path: tempPath, ThumbnailID: &thumb_id}
	db.Create(&insertFile)
	return c.JSON(http.StatusOK, insertFile)
}
