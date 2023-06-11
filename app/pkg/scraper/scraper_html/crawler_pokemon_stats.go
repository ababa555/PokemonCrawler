package scraper_html

import (
	"PokemonCrawler/app/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonStats ポケモンのステータスを取得します
func CreatePokemonStats(page *goquery.Document, pokemon models.Pokemon, version models.Version) models.PokemonStats {
	var hp, attack, defense, spAttack, spDefense, speed int
	if version.Name == "pika_vee" {
		hp, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(4) > div.table.layout_right > table > tbody > tr:nth-child(2) > td.left").Text()))
		attack, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(4) > div.table.layout_right > table > tbody > tr:nth-child(3) > td.left").Text()))
		defense, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(4) > div.table.layout_right > table > tbody > tr:nth-child(4) > td.left").Text()))
		spAttack, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(4) > div.table.layout_right > table > tbody > tr:nth-child(5) > td.left").Text()))
		spDefense, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(4) > div.table.layout_right > table > tbody > tr:nth-child(6) > td.left").Text()))
		speed, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(4) > div.table.layout_right > table > tbody > tr:nth-child(7) > td.left").Text()))
	} else {
		hp, _ = strconv.Atoi(removeRankingText(page.Find("#stats_anchor > table > tbody > tr:nth-child(2) > td.left").Text()))
		attack, _ = strconv.Atoi(removeRankingText(page.Find("#stats_anchor > table > tbody > tr:nth-child(3) > td.left").Text()))
		defense, _ = strconv.Atoi(removeRankingText(page.Find("#stats_anchor > table > tbody > tr:nth-child(4) > td.left").Text()))
		spAttack, _ = strconv.Atoi(removeRankingText(page.Find("#stats_anchor > table > tbody > tr:nth-child(5) > td.left").Text()))
		spDefense, _ = strconv.Atoi(removeRankingText(page.Find("#stats_anchor > table > tbody > tr:nth-child(6) > td.left").Text()))
		speed, _ = strconv.Atoi(removeRankingText(page.Find("#stats_anchor > table > tbody > tr:nth-child(7) > td.left").Text()))
	}

	pokemonStats := models.PokemonStats{
		PokemonID: pokemon.ID,
		Hp:        hp,
		Attack:    attack,
		Defense:   defense,
		SpAttack:  spAttack,
		SpDefense: spDefense,
		Speed:     speed,
	}

	return pokemonStats
}

func removeRankingText(target string) string {
	// 45(796位)のように取得されてくるので、順位の部分は削除する
	r := regexp.MustCompile(`\([0-9]*位\)`)
	value := r.ReplaceAllString(target, "")

	return removeNbsp(value)
}

func removeNbsp(target string) string {
	return strings.Replace(target, "\u00A0", "", 1)
}
