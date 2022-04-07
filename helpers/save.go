package helpers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// SaveImage 画像をダウンロードします
func SaveImage(imageURL string, filename string) {
	response, e := http.Get(imageURL)
	if e != nil {
		log.Fatal(e)
		log.Fatal(response)
	}

	exe, _ := os.Executable()
	current := filepath.Dir(exe)
	path := filepath.Join(current, "download", "image", filename)

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

// SaveCsv CSVをダウンロードします
func SaveCsv(data string, version string, filename string) {
	path := SavePath(version, filename)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0666)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// Openエラー処理
	}
	defer file.Close()

	file.Write(([]byte)(data))
}

func SavePath(version string, filename string) string {
	exe, _ := os.Executable()
	current := filepath.Dir(exe)
	path := filepath.Join(current, "download", "csv", version, filename)
	return path
}
