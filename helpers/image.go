package helpers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// ImageDownload 画像をダウンロードします
func ImageDownload(imageURL string, filename string) {
	response, e := http.Get(imageURL)
	if e != nil {
		log.Fatal(e)
		log.Fatal(response)
	}

	exe, _ := os.Executable()
	current := filepath.Dir(exe)
	path := filepath.Join(current, "images", filename)

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// file, err := os.Create(path.Join(current, filename))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = io.Copy(file, response.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
