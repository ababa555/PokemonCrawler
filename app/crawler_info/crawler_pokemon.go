package crawler_info

import (
	"PokemonCrawler/app/models"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func extractionWeightOrHeightFromString() func(target string) int {
	return func(target string) int {
		if target == "" {
			return 0
		}
		// 数字だけを取り出す
		r := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		numberStrings := r.FindAllString(target, -1)
		numberFloat, _ := strconv.ParseFloat(numberStrings[0], 32)
		numberInt, _ := strconv.Atoi(strconv.FormatFloat(numberFloat*10, 'g', 4, 64))
		return numberInt
	}
}

// CreatePokemon ポケモンの基本情報を取得します
func CreatePokemon(page *goquery.Document, id string, index int, isDefault bool, hasGaralNo bool, version models.Version) models.Pokemon {
	adjust := map[bool]int{true: 1, false: 0}[hasGaralNo]
	var no, heightAsString, weightAsString string
	if version.Name == "pika_vee" {
		no = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(4) > td:nth-child(2)").First().Text()
		heightAsString = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(5) > td:nth-child(2)").First().Text()
		weightAsString = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(6) > td:nth-child(2)").First().Text()
	} else {
		no = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(4+adjust) + ") > td:nth-child(2)").First().Text()
		heightAsString = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(5+adjust) + ") > td:nth-child(2)").First().Text()
		weightAsString = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(6+adjust) + ") > td:nth-child(2)").First().Text()
	}
	order := index

	e := extractionWeightOrHeightFromString()
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
