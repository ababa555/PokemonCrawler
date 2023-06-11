package models

import "github.com/PuerkitoBio/goquery"

type DocWithFile struct {
	Doc      *goquery.Document
	FileName string
	Version  Version
	Index    int
	SubIndex int
}
