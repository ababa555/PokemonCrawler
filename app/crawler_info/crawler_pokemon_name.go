package crawler_info

import (
	"PokemonCrawler/app/models"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// CreatePokemonNames ポケモンの名前を取得します（日本名＋英名）
func CreatePokemonName(page *goquery.Document, pokemon models.Pokemon, version models.Version) models.PokemonNames {
	var pokemonNames models.PokemonNames
	const jaLocalLanguageID = 1
	const enLocalLanguageID = 2

	if pokemon.ID[len(pokemon.ID)-4:len(pokemon.ID)-1] == "773" {
		typeAsEnum, _ := models.SilvallyTypeAsEnum(pokemon.ID[len(pokemon.ID)-1:])
		pokemonName := models.PokemonName{
			PokemonID:       pokemon.ID + "-" + strconv.Itoa(jaLocalLanguageID),
			LocalLanguageID: jaLocalLanguageID,
			Name:            "シルヴァディ",
			FormName:        models.TypeAsString(typeAsEnum),
		}
		pokemonNames = append(pokemonNames, pokemonName)

		pokemonNameEn := models.PokemonName{
			PokemonID:       pokemon.ID + "-" + strconv.Itoa(enLocalLanguageID),
			LocalLanguageID: enLocalLanguageID,
			Name:            "Silvally",
			FormName:        models.TypeAsString(typeAsEnum),
		}
		pokemonNames = append(pokemonNames, pokemonNameEn)
	} else {
		var name, nameEn string
		if version.Name == "pika_vee" {
			name = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr.head > th").First().Text()
		} else {
			name = page.Find("#base_anchor > table > tbody > tr.head > th").First().Text()
		}

		var formName = ""
		if !pokemon.IsDefault {
			page.Find(".select_list:not(.gen_list):first-child").Each(func(index int, s *goquery.Selection) {
				s.Find("li > strong").EachWithBreak(func(index int, s1 *goquery.Selection) bool {
					formName = s1.Text()
					return true
				})
			})
		}

		pokemonName := models.PokemonName{
			PokemonID:       pokemon.ID + "-" + strconv.Itoa(jaLocalLanguageID),
			LocalLanguageID: jaLocalLanguageID,
			Name:            name,
			FormName:        formName,
		}
		pokemonNames = append(pokemonNames, pokemonName)

		if version.Name == "pika_vee" {
			nameEn = page.Find("#contents > div:nth-child(4) > div.table.layout_left > table > tbody > tr:nth-child(8) > td:nth-child(2) > ul").First().Text()
		} else {
			adjust := 0
			for {
				nameEnText := page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(8+adjust) + ") > td.c1").First().Text()
				if nameEnText == "英語名" {
					break
				}
				adjust++
			}
			nameEn = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(8+adjust) + ") > td:nth-child(2) > ul > li").First().Text()
		}

		pokemonNameEn := models.PokemonName{
			PokemonID:       pokemon.ID + "-" + strconv.Itoa(enLocalLanguageID),
			LocalLanguageID: enLocalLanguageID,
			Name:            nameEn,
			FormName:        formName,
		}
		pokemonNames = append(pokemonNames, pokemonNameEn)
	}

	return pokemonNames
}
