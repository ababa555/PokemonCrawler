package crawler_info

import (
	"PokemonCrawler/app/models"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonAbility ポケモンの特性を取得します
func CreatePokemonAbility(page *goquery.Document, pokemon models.Pokemon) models.PokemonAbilities {
	index := 34
	adjust := 0
	for {
		text := page.Find("#stats_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index+adjust) + ") > th").Text()
		if !strings.Contains(text, "特性(とくせい)") {
			adjust++
			continue
		}
		break
	}
	// 特性
	var pokemonAbilities models.PokemonAbilities
	for {
		text := page.Find("#stats_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index+1+adjust) + ") > td.c1 > a").Text()
		if text == "" {
			break
		}
		pokemonAbility := models.PokemonAbility{
			PokemonID:   pokemon.ID,
			AbilityName: text,
			IsHidden:    false,
		}
		pokemonAbilities = append(pokemonAbilities, pokemonAbility)
		index++
	}
	// 隠れ特性
	for {
		text := page.Find("#stats_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(index+3+adjust) + ") > td.c1 > a").Text()
		if text == "" {
			break
		}
		pokemonAbility := models.PokemonAbility{
			PokemonID:   pokemon.ID,
			AbilityName: text,
			IsHidden:    true,
		}
		pokemonAbilities = append(pokemonAbilities, pokemonAbility)
		index++
	}
	return pokemonAbilities
}
