package crawler_info

import (
	"PokemonCrawler/app/models"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonEvolutionChains ポケモンの進化情報を取得します
func CreatePokemonEvolutionChains(page *goquery.Document, pokemon models.Pokemon, version models.Version) models.PokemonEvolutionChains {
	var pokemonEvolutionChains models.PokemonEvolutionChains

	evolutions := page.Find(".evo_list > li > a")
	if evolutions.Nodes == nil {
		pokemonEvolutionChain := models.PokemonEvolutionChain{
			PokemonID:        pokemon.ID,
			EvolutionChainID: pokemon.ID,
			Order:            1,
		}
		pokemonEvolutionChains = append(pokemonEvolutionChains, pokemonEvolutionChain)
		return pokemonEvolutionChains
	}

	var addEvolutionChainIDList []string
	evolutions.Each(func(index int, s *goquery.Selection) {
		text, _ := s.Attr("href")
		slice := strings.Split(text, "/")
		evolutionChainID := slice[1]

		r := regexp.MustCompile(`n[0-9]*`)
		if !r.MatchString(evolutionChainID) {
			return
		}

		if contains(addEvolutionChainIDList, evolutionChainID) {
			return
		}

		pokemonEvolutionChain := models.PokemonEvolutionChain{
			PokemonID:        pokemon.ID,
			EvolutionChainID: version.No + "-" + evolutionChainID,
			Order:            index + 1,
		}
		pokemonEvolutionChains = append(pokemonEvolutionChains, pokemonEvolutionChain)
		addEvolutionChainIDList = append(addEvolutionChainIDList, evolutionChainID)
	})

	return pokemonEvolutionChains
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
