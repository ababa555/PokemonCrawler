package logics

import (
	"bufio"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func visit(searchurl string) *goquery.Document {
	resp, err := http.Get(searchurl)
	if err != nil {
		return nil
	}

	// EUC_JPからUTF8に変換
	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.EUCJP.NewDecoder())
	doc, _ := goquery.NewDocumentFromReader(utfBody)

	return doc
}
