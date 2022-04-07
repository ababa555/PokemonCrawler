package crawler_info

import (
	"PokemonCrawler/app/models"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonMove ポケモンの覚える技を取得します
func CreatePokemonMove(page *goquery.Document, pokemon models.Pokemon) models.PokemonMoves {
	var pokemonMoves models.PokemonMoves
	//page.Find("#move_anchor").Each(func(index int, s *goquery.Selection) {
	page.Find("#move_list > tbody > tr").EachWithBreak(func(index int, s1 *goquery.Selection) bool {
		name := s1.Find("td.move_name_cell > a").Text()
		title := s1.Find("#move_list > tbody > tr:nth-child(" + strconv.Itoa(index+1) + ") > th").Text()
		if strings.Contains(title, "過去作でしか覚えられない技") {
			return false
		}
		if name != "" {
			pokemonMove := models.PokemonMove{
				PokemonID: pokemon.ID,
				MoveName:  name,
			}
			pokemonMoves = append(pokemonMoves, pokemonMove)
		}
		return true
	})
	//})
	return pokemonMoves
}
