package scraper_html

import (
	"PokemonCrawler/app/models"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonMove ポケモンの覚える技を取得します
func CreatePokemonMove(page *goquery.Document, pokemon models.Pokemon) models.PokemonMoves {
	var pokemonMoves models.PokemonMoves

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
			if existInPokemonMoves(pokemonMoves, pokemonMove) {
				return true
			}
			pokemonMoves = append(pokemonMoves, pokemonMove)
		}
		return true
	})

	return pokemonMoves
}

func existInPokemonMoves(s []models.PokemonMove, e models.PokemonMove) bool {
	for _, a := range s {
		if ok := reflect.DeepEqual(a, e); ok {
			return true
		}
	}
	return false
}
