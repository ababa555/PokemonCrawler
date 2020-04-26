package logics

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func createSearchURL(version string, no string) string {
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

func createSearchURLMovePage(version string, urlParams string) string {
	baseURL, _ := url.Parse(strings.Join([]string{
		"https://yakkun.com/",
		version,
		"/",
		"move_list.htm",
		urlParams,
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
	searchurl := createSearchURL(version, searchNo)
	page, statusCode := visitImpl(searchurl, checkRedirect)
	if page == nil {
		return nil, statusCode
	}
	return page, statusCode
}

// VisitMovePage 技のページを読み込みます
func VisitMovePage(version string, searchNo string) *goquery.Document {
	searchurl := createSearchURLMovePage(version, searchNo)
	page, _ := visitImpl(searchurl, false)
	if page == nil {
		panic("page not found")
	}
	return page
}
