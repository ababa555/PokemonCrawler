package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"./logics"
	"./models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
)

func savePath(version string, filename string) string {
	exe, _ := os.Executable()
	current := filepath.Dir(exe)
	path := filepath.Join(current, "files", version, filename)
	return path
}

func saveCsv(data string, version string, filename string) {
	path := savePath(version, filename)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0666)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// Openエラー処理
	}
	defer file.Close()

	file.Write(([]byte)(data))
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

func createMovesCsv(version string) {
	path := savePath(version, "moves.csv")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		var result = "\"id\",\"version\",\"name\",\"type\",\"power1\",\"power2\",\"pp\",\"accuracy\",\"priority\",\"damageType\",\"isDirect\",\"canProtect\"" + "\n"
		moves := logics.CretaMoves(version)
		rt := reflect.New(reflect.TypeOf(models.Move{})).Elem().Type()
		for _, move := range moves {
			result += createData(move, rt)
		}
		saveCsv(result, version, "moves.csv")
	}
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

func createPokemonInfoCsv(page *goquery.Document, id string, index int, isDefault bool, hasGaralNo bool, version string) {
	pokemon := logics.CreatePokemon(page, id, index, isDefault, hasGaralNo, version)
	pokemonNames := logics.CreatePokemonNames(page, pokemon, hasGaralNo, version)
	pokemonStats := logics.CreatePokemonStats(page, pokemon, version)
	pokemonTypes := logics.CreatePokemonTypes(page, pokemon)
	pokemonMoves := logics.CreatePokemonMoves(page, pokemon, version)
	pokemonEvolutionChains := logics.CreatePokemonEvolutionChains(page, pokemon)

	createPokemonCsv(pokemon, version, id)
	createPokemonNamesCsv(pokemonNames, version, id)
	createPokemonStatsCsv(pokemonStats, version, id)
	createPokemonTypesCsv(pokemonTypes, version, id)
	createPokemonMovesCsv(pokemonMoves, version, id)
	createPokemonEvolutionChainsCsv(pokemonEvolutionChains, version, id)

	// ピカブイは特性がない
	if version != "pika_vee" {
		pokemonAbilities := logics.CreatePokemonAbilities(page, pokemon)
		createPokemonAbilitiesCsv(pokemonAbilities, version, id)
	}
}

func createPokemonNamesCsv(data models.PokemonNames, version string, id string) {
	const fileName = "pokemonNames.csv"
	var list models.PokemonNames
	list = append(list, data...)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"localLanguageId\",\"name\",\"formName\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonName{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func createPokemonStatsCsv(data models.PokemonStats, version string, id string) {
	const fileName = "pokemonStatses.csv"
	var list models.PokemonStatses
	list = append(list, data)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"hp\",\"attack\",\"defense\",\"spAttack\",\"spDefense\",\"speed\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonStats{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func createPokemonTypesCsv(data models.PokemonTypes, version string, id string) {
	const fileName = "pokemonTypes.csv"
	var list models.PokemonTypes
	list = append(list, data...)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"typeId\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonType{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func createPokemonAbilitiesCsv(data models.PokemonAbilities, version string, id string) {
	const fileName = "pokemonAbilities.csv"
	var list models.PokemonAbilities
	list = append(list, data...)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"abilityName\",\"isHidden\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonAbility{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func createPokemonMovesCsv(data models.PokemonMoves, version string, id string) {
	const fileName = "pokemonMoves.csv"
	var list models.PokemonMoves
	list = append(list, data...)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"moveName\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonMove{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func createPokemonEvolutionChainsCsv(data models.PokemonEvolutionChains, version string, id string) {
	const fileName = "pokemonEvolutionChains.csv"
	var list models.PokemonEvolutionChains
	list = append(list, data...)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"pokemonId\",\"evolutionChainId\",\"order\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.PokemonEvolutionChain{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func createPokemonCsv(data models.Pokemon, version string, id string) {
	const fileName = "pokemons.csv"
	var list models.Pokemons
	list = append(list, data)
	path := savePath(version, fileName)
	var result string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result = "\"id\",\"no\",\"height\",\"weight\",\"order\",\"isDefault\"" + "\n"
	}

	rt := reflect.New(reflect.TypeOf(models.Pokemon{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	saveCsv(result, version, fileName)
}

func isDownloaded(version string, id string) bool {
	path := savePath(version, "pokemons.csv")
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

func main() {
	println("処理を開始します...")

	versions := [3]string{
		"swsh",
		"pika_vee",
		"sm",
	}

	// 技一覧
	for _, version := range versions {
		createMovesCsv(version)
	}

	for _, version := range versions {
		index := 1
		for {
			// ここに無限ループで実行する処理を記述
			searchNo := "n" + strconv.Itoa(index)
			println("処理中... " + searchNo + "(" + version + ")")

			// ダウンロード済の場合、次のポケモンへ
			if isDownloaded(version, searchNo) {
				index++
				continue
			}

			page, statusCode := logics.Visit(version, searchNo, true)
			time.Sleep(5000 * time.Millisecond)
			if statusCode == 404 {
				break
			} else if page == nil {
				index++
				continue
			}

			// ガラル地方のポケモンかどうかを調べる
			hasGaralNo := getHasGaralNo(version, page)

			// ポケモンの情報をcsvに保存
			createPokemonInfoCsv(page, searchNo, index, true, hasGaralNo, version)

			// 他のフォルム
			forms := getForms(page)
			for _, form := range forms {
				searchNo := "n" + form
				println("処理中...:" + searchNo + "(" + version + ")")

				page, _ := logics.Visit(version, searchNo, true)
				time.Sleep(5000 * time.Millisecond)

				if page != nil {
					// ガラル地方のポケモンかどうかを調べる
					hasGaralNo := getHasGaralNo(version, page)

					// ポケモンの情報をcsvに保存
					createPokemonInfoCsv(page, searchNo, index, false, hasGaralNo, version)
				}
			}

			index++
		}
	}
}
