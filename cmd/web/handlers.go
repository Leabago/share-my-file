package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"share-my-file/pkg/forms"
	"share-my-file/pkg/models"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(":id")
	fullURL := r.Host + r.URL.Path

	var png []byte
	png, err := qrcode.Encode(fullURL, qrcode.Medium, 256)
	if err != nil {
		app.logger.infoLog.Println("cant create QRcode")
	}
	base64ImageData := base64.StdEncoding.EncodeToString(png)

	fileNameList := app.redisClient.LRange(app.getAvailableKey(code), 0, -1).Val()
	fmt.Println("fileNames:", fileNameList)

	var file = &models.File{
		FolderCode:   code,
		Exist:        app.fileExist(code),
		URL:          fullURL,
		QRcodeBase64: base64ImageData,
		FileNameList: fileNameList,
	}

	// s, err := app.files.Get()

	// if err != nil {
	// 	if errors.Is(err, models.ErrNoRecord) {
	// 		app.notFound(w)
	// 	} else {
	// 		app.serverError(w, err)
	// 	}
	// 	return
	// }

	app.render(w, r, "show.page.tmpl.html", &templateData{
		File: file,
	})

	// filePath := folderPath + folderBegin + code + zipName
	// filename := folderBegin + code + zipName

	// w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	// w.Header().Set("Content-Type", "application/octet-stream")
	// http.ServeFile(w, r, filePath)

}

func (app *application) getSnippet(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(":id")

	// var file = &models.File{
	// 	FolderCode: code,
	// }

	// app.render(w, r, "show.page.tmpl.html", &templateData{
	// 	File: file,
	// })

	if app.fileExist(code) {

		filename := folderBegin + code + zipName
		filePath := folderPath + filename

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, filePath)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/archive/%s", code), http.StatusSeeOther)
	}
}

// homeGetFiles upload files to zip
func (app *application) homeGetFiles(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(":code")
	// get folder name

	fmt.Print("code:= ", code)
	var zipFileName = folderPath + folderBegin + code + zipName
	app.logger.infoLog.Printf("create new folder %s", zipFileName)
	fileNameList, err := ParseMediaType(r, zipFileName, app.maxFileSize)
	if err != nil {
		app.serverError(w, err)
	}

	// send key to redis. key is going to expire
	app.redisClient.RPush((app.getAvailableKey(code)), fileNameList)
	app.redisClient.Expire(app.getAvailableKey(code), smallTime).Result()
	// w.Write([]byte(code))

	http.Redirect(w, r, ("/archive/" + code), http.StatusSeeOther)
}

func (app *application) getAvailableKey(code string) string {
	return available + code
}

func (app *application) redirectToArchive(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	title := form.Values.Get("title")
	app.logger.infoLog.Printf("redirectToArchive %s", title)

}

func (app *application) createDownloadForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "download.page.tmpl.html", &templateData{})
}

func (application *application) redirectHome(w http.ResponseWriter, r *http.Request) {
	code := createUserCode()
	http.Redirect(w, r, "/upload/"+code, http.StatusSeeOther)
}

func check(e error) {
	if e != nil {
		fmt.Println("panic panic")
		panic(e)
	}
}
