package crawler_info

import (
	"PokemonCrawler/app/models"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonType ポケモンのタイプを取得します
func CreatePokemonType(page *goquery.Document, pokemon models.Pokemon) models.PokemonTypes {
	var pokemonTypes models.PokemonTypes
	s := page.Find(".type").First()
	s.Find("li > a > img").Each(func(_ int, s2 *goquery.Selection) {
		typeAsString, _ := s2.Attr("alt")
		typeID, _ := models.TypeAsEnum(typeAsString)

		pokemonType := models.PokemonType{
			PokemonID: pokemon.ID,
			TypeID:    typeID,
		}
		pokemonTypes = append(pokemonTypes, pokemonType)
	})
	return pokemonTypes
}
