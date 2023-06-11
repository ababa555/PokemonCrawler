package helpers

import (
	. "PokemonCrawler/app/models"
	"bufio"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func CreateSearchURL(version string, no string) string {
	baseURL, _ := url.Parse(strings.Join([]string{
		"https://yakkun.com/",
		version,
		"/",
		"zukan",
		"/",
		no,
	}, ""))

	return baseURL.String()
}

func CreateSearchURLMovePage(version string, urlParams string) string {
	baseURL, _ := url.Parse(strings.Join([]string{
		"https://yakkun.com/",
		version,
		"/",
		"move_list.htm",
		urlParams,
	}, ""))

	return baseURL.String()
}

func CreateSearchURLZukan(no string) string {
	baseURL, _ := url.Parse(strings.Join([]string{
		"https://zukan.pokemon.co.jp/detail/",
		no,
	}, ""))

	return baseURL.String()
}

func visitImpl(searchurl string, checkRedirect bool) (*goquery.Document, int) {
	var resp *http.Response
	var err error
	if checkRedirect {
		var RedirectAttemptedError = errors.New("redirect")
		client := &http.Client{
			Timeout: time.Duration(3) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return RedirectAttemptedError
			},
		}

		resp, err = client.Get(searchurl)
		if err != nil {
			// ソードシールドからサンムーンなどゲームバージョンが異なるページに飛ばされた場合はnilを返す
			expect := strings.Split(searchurl, "/")[3]
			r := reflect.ValueOf(err)
			redirectURL := reflect.Indirect(r).FieldByName("URL").String()
			actually := strings.Split(redirectURL, "/")[1]
			if expect != actually {
				return nil, resp.StatusCode
			}
		}
	} else {
		resp, err = http.Get(searchurl)
		if err != nil {
			return nil, resp.StatusCode
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, resp.StatusCode
	}

	// EUC_JPからUTF8に変換
	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.EUCJP.NewDecoder())
	doc, _ := goquery.NewDocumentFromReader(utfBody)

	return doc, resp.StatusCode
}

// Visit ページを読み込みます
func Visit(version string, searchNo string, checkRedirect bool) (*goquery.Document, int) {
	searchurl := CreateSearchURL(version, searchNo)
	page, statusCode := visitImpl(searchurl, checkRedirect)
	if page == nil {
		return nil, statusCode
	}
	return page, statusCode
}

// VisitHtml ページを読み込みます
func VisitHtml(version Version) []DocWithFile {
	fileName := SavePath("html\\pokemon", version.Id())
	files, err := os.ReadDir(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var docs []DocWithFile
	subIndex := 1
	preNo := ""
	for _, file := range files {
		file, err := os.Open(fileName + "\\" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			log.Fatal("Error loading HTML into goquery:", err)
		}

		filenameWithExt := filepath.Base(file.Name())
		filename := strings.TrimSuffix(filenameWithExt, filepath.Ext(filenameWithExt))

		re := regexp.MustCompile(`\d+`)
		index, _ := strconv.Atoi(re.FindString(filename))

		reNto := regexp.MustCompile(`n\d+`)
		no := reNto.FindString(filename)
		if preNo != no {
			subIndex = 1
		} else {
			subIndex++
		}

		docs = append(docs, DocWithFile{Doc: doc, FileName: filename, Version: version, Index: index, SubIndex: subIndex})
	}
	return docs
}

// VisitMovePage 技のページを読み込みます
func VisitMovePage(version string, searchNo string) *goquery.Document {
	searchurl := CreateSearchURLMovePage(version, searchNo)
	page, _ := visitImpl(searchurl, false)
	if page == nil {
		panic("page not found")
	}
	return page
}

// VisitMoveHtml　技のhtmlを読み込みます
func VisitMoveHtml(version string, params string) *goquery.Document {
	fileName := SavePath("html\\move", version, params, "move.html")
	file, err := os.Open(fileName) // ファイル名を適切に変更してください
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal("Error loading HTML into goquery:", err)
	}

	return doc
}

// VisitZukan ポケモン図鑑のページを読み込みます
func VisitZukan(searchNo string) (*goquery.Document, int) {
	searchurl := CreateSearchURLZukan(searchNo)
	page, statusCode := visitImpl(searchurl, false)
	if page == nil {
		return nil, statusCode
	}
	return page, statusCode
}
