package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	fileName := handler.Filename
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, file)
	json.NewEncoder(w).Encode(map[string]string{
		"image_url": "http://localhost:10000/images/" + fileName,
	})
}

func UploadCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	fileName := "import.csv"
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, file)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "upload csv " + fileName + " success",
	})
}
