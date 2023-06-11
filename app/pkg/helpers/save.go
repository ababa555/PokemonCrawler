package helpers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// SaveImage 画像をダウンロードします
func SaveImage(imageURL string, filename string) error {
	response, e := http.Get(imageURL)
	if e != nil {
		log.Fatal(e)
		log.Fatal(response)
	}
	defer response.Body.Close()

	path := SavePath("image", "", filename)

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// SaveCsv CSVをダウンロードします
func SaveCsv(data string, version string, filename string) error {
	path := SavePath("csv", version, filename)

	err := EnsureDirExists(path)
	if err != nil {
		return err
	}

	err = SaveToFile(path, data)
	if err != nil {
		return err
	}

	return nil
}

func SavePath(category string, version string, subdirs ...string) string {
	current := CurrentDir()
	pathSegments := append([]string{current, "download", category, version}, subdirs...)
	path := filepath.Join(pathSegments...)
	return path
}

func CurrentDir() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(fmt.Errorf("os.Executable failed: %v", err))
	}
	workingDir := filepath.Dir(exe)

	for {
		// PokemonCrawlerが見つかったらここがルートディレクトリ
		if filepath.Base(workingDir) == "PokemonCrawler" {
			return workingDir
		}

		next := filepath.Dir(workingDir)
		if next == workingDir {
			break
		}

		workingDir = next
	}

	return ""
}

func EnsureDirExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("os.MkdirAll failed: %v", err)
		}
	}
	return nil
}

func SaveToFile(path string, data string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile failed: %v", err)
	}
	defer file.Close()

	_, err = file.Write(([]byte)(data))
	if err != nil {
		return fmt.Errorf("file.Write failed: %v", err)
	}

	return nil
}
