package scraper_html

import (
	"PokemonCrawler/app/models"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func extractionIntFromString() func(target string) string {
	return func(target string) string {
		if target == "" {
			return "0"
		}
		// 数字だけを取り出す
		r := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		numberStrings := r.FindAllString(target, -1)
		// numberFloat, _ := strconv.ParseFloat(numberStrings[0], 32)
		// numberInt, _ := strconv.Atoi(strconv.FormatFloat(numberFloat*10, 'g', 4, 64))
		return numberStrings[0]
	}
}

// CreatePokemon ポケモンの基本情報を取得します
func CreatePokemon(page *goquery.Document, id string, index int, subIndex int, isDefault bool, version models.Version) models.Pokemon {
	var no, heightAsString, weightAsString string
	if version.Name == "pika_vee" {
		no = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(4) > td:nth-child(2)").First().Text()
		heightAsString = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(5) > td:nth-child(2)").First().Text()
		weightAsString = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(6) > td:nth-child(2)").First().Text()
	} else {
		for i := 0; ; i++ {
			index := 4 + i
			text := page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index) + ") > td.c1").First().Text()
			if text == "全国No." || text == "ぜんこくNo." {
				no = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index) + ") > td:nth-child(2)").First().Text()
			}

			if text == "高さ" {
				heightAsString = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index) + ") > td:nth-child(2)").First().Text()
			}

			if text == "重さ" {
				weightAsString = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index) + ") > td:nth-child(2) > ul > li:nth-child(1)").First().Text()
			}

			if no != "" && heightAsString != "" && weightAsString != "" {
				break
			}
		}
	}
	order := subIndex

	e := extractionIntFromString()
	height := e(heightAsString)
	weight := e(weightAsString)

	pokemon := models.Pokemon{
		ID:        version.No + "-" + id,
		No:        no,
		Height:    height,
		Weight:    weight,
		Order:     order,
		IsDefault: isDefault,
	}

	return pokemon
}
