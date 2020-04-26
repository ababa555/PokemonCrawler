package logics

import (
	"regexp"
	"strconv"
	"strings"

	"../models"
	"github.com/PuerkitoBio/goquery"
)

func extractionWeightOrHeightFromString() func(target string) int {
	return func(target string) int {
		if target == "" {
			return 0
		}
		// 数字だけを取り出す
		pattern := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		numberStrings := pattern.FindAllString(target, -1)
		numberFloat, _ := strconv.ParseFloat(numberStrings[0], 32)
		numberInt, _ := strconv.Atoi(strconv.FormatFloat(numberFloat*10, 'g', 4, 64))
		return numberInt
	}
}

// CreatePokemon ポケモンの基本情報を取得します
func CreatePokemon(page *goquery.Document, id string, index int, isDefault bool, hasGaralNo bool, version string) models.Pokemon {
	adjust := map[bool]int{true: 1, false: 0}[hasGaralNo]
	var no, heightAsString, weightAsString string
	if version == "pika_vee" {
		no = page.Find("#contents > div:nth-child(5) > div.table.layout_left > table > tbody > tr:nth-child(4) > td:nth-child(2)").First().Text()
		heightAsString = page.Find("#contents > div:nth-child(5) > div.table.layout_left > table > tbody > tr:nth-child(5) > td:nth-child(2)").First().Text()
		weightAsString = page.Find("#contents > div:nth-child(5) > div.table.layout_left > table > tbody > tr:nth-child(6) > td:nth-child(2)").First().Text()
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
		ID:        id,
		No:        no,
		Height:    height,
		Weight:    weight,
		Order:     order,
		IsDefault: isDefault,
	}

	return pokemon
}

// CreatePokemonNames ポケモンの名前を取得します（日本名＋英名）
func CreatePokemonNames(page *goquery.Document, pokemon models.Pokemon, hasGaralNo bool, version string) models.PokemonNames {
	var pokemonNames models.PokemonNames
	var name, nameEn string
	if version == "pika_vee" {
		name = page.Find("#contents > div:nth-child(5) > div.table.layout_left > table > tbody > tr.head > th").First().Text()
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
		PokemonID:       pokemon.ID,
		LocalLanguageID: 1,
		Name:            name,
		FormName:        formName,
	}
	pokemonNames = append(pokemonNames, pokemonName)

	adjust := map[bool]int{true: 1, false: 0}[hasGaralNo]
	if version == "pika_vee" {
		nameEn = page.Find("#contents > div:nth-child(5) > div.table.layout_left > table > tbody > tr:nth-child(8) > td:nth-child(2) > ul").First().Text()
	} else {
		nameEn = page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(8+adjust) + ") > td:nth-child(2) > ul > li").First().Text()
	}
	pokemonNameEn := models.PokemonName{
		PokemonID:       pokemon.ID,
		LocalLanguageID: 2,
		Name:            nameEn,
		FormName:        formName,
	}
	pokemonNames = append(pokemonNames, pokemonNameEn)

	return pokemonNames
}

// CreatePokemonStats ポケモンのステータスを取得します
func CreatePokemonStats(page *goquery.Document, pokemon models.Pokemon, version string) models.PokemonStats {
	var hp, attack, defense, spAttack, spDefense, speed int
	if version == "pika_vee" {
		hp, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(5) > div.table.layout_right > table > tbody > tr:nth-child(2) > td.left").Text()))
		attack, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(5) > div.table.layout_right > table > tbody > tr:nth-child(3) > td.left").Text()))
		defense, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(5) > div.table.layout_right > table > tbody > tr:nth-child(4) > td.left").Text()))
		spAttack, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(5) > div.table.layout_right > table > tbody > tr:nth-child(5) > td.left").Text()))
		spDefense, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(5) > div.table.layout_right > table > tbody > tr:nth-child(6) > td.left").Text()))
		speed, _ = strconv.Atoi(removeNbsp(page.Find("#contents > div:nth-child(5) > div.table.layout_right > table > tbody > tr:nth-child(7) > td.left").Text()))
	} else {
		hp, _ = strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(2) > td.left").Text()))
		attack, _ = strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(3) > td.left").Text()))
		defense, _ = strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(4) > td.left").Text()))
		spAttack, _ = strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(5) > td.left").Text()))
		spDefense, _ = strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(6) > td.left").Text()))
		speed, _ = strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(7) > td.left").Text()))
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

// CreatePokemonTypes ポケモンのタイプを取得します
func CreatePokemonTypes(page *goquery.Document, pokemon models.Pokemon) models.PokemonTypes {
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

// CreatePokemonAbilities ポケモンの特性を取得します
func CreatePokemonAbilities(page *goquery.Document, pokemon models.Pokemon) models.PokemonAbilities {
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

// CreatePokemonMoves ポケモンの覚える技を取得します
func CreatePokemonMoves(page *goquery.Document, pokemon models.Pokemon, version string) models.PokemonMoves {
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

// CreatePokemonEvolutionChains ポケモンの進化情報を取得します
func CreatePokemonEvolutionChains(page *goquery.Document, pokemon models.Pokemon) models.PokemonEvolutionChains {
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

	evolutions.Each(func(index int, s *goquery.Selection) {
		text, _ := s.Attr("href")
		slice := strings.Split(text, "/")
		evolutionChainID := slice[1]
		pokemonEvolutionChain := models.PokemonEvolutionChain{
			PokemonID:        pokemon.ID,
			EvolutionChainID: evolutionChainID,
			Order:            index + 1,
		}
		pokemonEvolutionChains = append(pokemonEvolutionChains, pokemonEvolutionChain)
	})

	return pokemonEvolutionChains
}

func removeNbsp(target string) string {
	return strings.Replace(target, "\u00A0", "", 1)
}
