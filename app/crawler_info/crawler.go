package crawler_info

import (
	"PokemonCrawler/app/models"
	"PokemonCrawler/helpers"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
)

func Run() {
	versions := []models.Version{
		*models.NewVersion("1", "sm"),
		*models.NewVersion("2", "pika_vee"),
		*models.NewVersion("3", "swsh"),
	}

	// // 技一覧
	// for _, version := range versions {
	// 	createMovesCsv(version)
	// }

	for _, version := range versions {
		index := 1
		for {
			// ここに無限ループで実行する処理を記述
			searchNo := "n" + strconv.Itoa(index)
			println("処理中... " + searchNo + "(" + version.Name + ")")

			// ダウンロード済の場合、次のポケモンへ
			if isDownloaded(version.Id(), searchNo) {
				index++
				continue
			}

			page, statusCode := helpers.Visit(version.Name, searchNo, true)
			time.Sleep(2000 * time.Millisecond)
			if statusCode == 404 {
				break
			} else if page == nil {
				index++
				continue
			}

			// ガラル地方のポケモンかどうかを調べる
			hasGaralNo := getHasGaralNo(version.Name, page)

			// ポケモンの情報をcsvに保存
			createPokemonInfoCsv(page, searchNo, index, true, hasGaralNo, version)

			// 他のフォルム
			forms := getForms(page)
			for _, form := range forms {
				searchNo := "n" + form
				println("処理中...:" + searchNo + "(" + version.Name + ")")

				page, _ := helpers.Visit(version.Name, searchNo, true)
				time.Sleep(5000 * time.Millisecond)

				if page != nil {
					// ガラル地方のポケモンかどうかを調べる
					hasGaralNo := getHasGaralNo(version.Name, page)

					// ポケモンの情報をcsvに保存
					createPokemonInfoCsv(page, searchNo, index, false, hasGaralNo, version)
				}
			}

			index++
		}
	}
}

func createMovesCsv(version models.Version) {
	path := helpers.SavePath(version.Id(), "moves.csv")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		var result = "\"id\",\"name\",\"type\",\"power1\",\"power2\",\"pp\",\"accuracy\",\"priority\",\"damageType\",\"isDirect\",\"canProtect\"" + "\n"
		moves := CretaMoves(version)
		rt := reflect.New(reflect.TypeOf(models.Move{})).Elem().Type()
		for _, move := range moves {
			result += createData(move, rt)
		}
		helpers.SaveCsv(result, version.Id(), "moves.csv")
	}
}

func createData(data interface{}, rt reflect.Type) string {
	var result string
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		r := reflect.ValueOf(data)
		var value string
		if f.Type.String() == "string" {
			value = reflect.Indirect(r).FieldByName(f.Name).Interface().(string)
		} else if f.Type.String() == "int" {
			value = strconv.Itoa(reflect.Indirect(r).FieldByName(f.Name).Interface().(int))
		} else if f.Type.String() == "models.Type" {
			t := reflect.Indirect(r).FieldByName(f.Name).Interface().(models.Type)
			value = strconv.Itoa(int(t))
		} else if f.Type.String() == "bool" {
			b := reflect.Indirect(r).FieldByName(f.Name).Interface().(bool)
			value = strconv.FormatBool(b)
		} else {
			value = reflect.Indirect(r).FieldByName(f.Name).Interface().(string)
		}

		value = "\"" + value + "\""
		if i == 0 {
			result = value
		} else {
			result = result + "," + value
		}
	}
	return result + "\n"
}

func getHasGaralNo(version string, page *goquery.Document) bool {
	if version == "swsh" {
		text := page.Find("#base_anchor > table > tbody > tr:nth-child(4) > td.c1").Text()
		if text == "ガラルNo." {
			return true
		}
	}
	return false
}

func createPokemonInfoCsv(page *goquery.Document, id string, index int, isDefault bool, hasGaralNo bool, version models.Version) {
	pokemon := CreatePokemon(page, id, index, isDefault, hasGaralNo, version)
	pokemonNames := CreatePokemonName(page, pokemon, hasGaralNo, version)
	pokemonStats := CreatePokemonStats(page, pokemon, version)
	pokemonTypes := CreatePokemonType(page, pokemon)
	pokemonMoves := CreatePokemonMove(page, pokemon)
	pokemonEvolutionChains := CreatePokemonEvolutionChains(page, pokemon, version)

	createPokemonCsv(pokemon, version.Id(), id)
	createPokemonNamesCsv(pokemonNames, version.Id(), id)
	createPokemonStatsCsv(pokemonStats, version.Id(), id)
	createPokemonTypesCsv(pokemonTypes, version.Id(), id)
	createPokemonMovesCsv(pokemonMoves, version.Id(), id)
	createPokemonEvolutionChainsCsv(pokemonEvolutionChains, version.Id(), id)

	// ピカブイは特性がない
	if version.Name != "pika_vee" {
		pokemonAbilities := CreatePokemonAbility(page, pokemon)
		createPokemonAbilitiesCsv(pokemonAbilities, version.Id(), id)
	}
}

func createPokemonNamesCsv(data models.PokemonNames, version string, id string) {
	const fileName = "pokemonNames.csv"
	var list models.PokemonNames
	list = append(list, data...)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"localLanguageId\",\"name\",\"formName\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonName{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func createPokemonStatsCsv(data models.PokemonStats, version string, id string) {
	const fileName = "pokemonStatses.csv"
	var list models.PokemonStatses
	list = append(list, data)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"hp\",\"attack\",\"defense\",\"spAttack\",\"spDefense\",\"speed\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonStats{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func createPokemonTypesCsv(data models.PokemonTypes, version string, id string) {
	const fileName = "pokemonTypes.csv"
	var list models.PokemonTypes
	list = append(list, data...)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"typeId\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonType{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func createPokemonAbilitiesCsv(data models.PokemonAbilities, version string, id string) {
	const fileName = "pokemonAbilities.csv"
	var list models.PokemonAbilities
	list = append(list, data...)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"abilityName\",\"isHidden\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonAbility{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func createPokemonMovesCsv(data models.PokemonMoves, version string, id string) {
	const fileName = "pokemonMoves.csv"
	var list models.PokemonMoves
	list = append(list, data...)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"moveName\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonMove{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func createPokemonEvolutionChainsCsv(data models.PokemonEvolutionChains, version string, id string) {
	const fileName = "pokemonEvolutionChains.csv"
	var list models.PokemonEvolutionChains
	list = append(list, data...)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"evolutionChainId\",\"order\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonEvolutionChain{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func createPokemonCsv(data models.Pokemon, version string, id string) {
	const fileName = "pokemons.csv"
	var list models.Pokemons
	list = append(list, data)
	path := helpers.SavePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"id\",\"no\",\"height\",\"weight\",\"order\",\"isDefault\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.Pokemon{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func getForms(page *goquery.Document) []string {
	var forms []string
	// 他のフォルムがあるか
	page.Find(".select_list:not(.gen_list):first-child").Each(func(index int, s *goquery.Selection) {
		s.Find("li > a").Each(func(index int, s1 *goquery.Selection) {
			text, _ := s1.Attr("href")
			slice := strings.Split(text, "/")
			form := strings.Replace(slice[3], "n", "", 1)
			forms = append(forms, form)
		})
	})
	return forms
}

func isDownloaded(version string, id string) bool {
	path := helpers.SavePath(version, "pokemons.csv")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		csv, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer csv.Close()

		structs := []*models.Pokemon{}

		if err := gocsv.UnmarshalFile(csv, &structs); err != nil { // Load clients from file
			panic(err)
		}

		for _, row := range structs {
			if row.ID == id {
				return true
			}
		}
	}
	return false
}
