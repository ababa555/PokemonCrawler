package scraper_html

import (
	"PokemonCrawler/app/models"
	. "PokemonCrawler/app/models"
	"PokemonCrawler/app/pkg/helpers"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
)

func Run() {
	versions := []models.Version{
		*models.NewVersion("1", "sm"),
		*models.NewVersion("2", "pika_vee"),
		*models.NewVersion("3", "swsh"),
		*models.NewVersion("4", "sv"),
	}

	// 技一覧
	for _, version := range versions {
		createMovesCsv(version)
	}

	var docWithFiles []DocWithFile
	for _, version := range versions {
		docs := helpers.VisitHtml(version)
		docWithFiles = append(docWithFiles, docs...)
	}

	for _, docWithFile := range docWithFiles {
		println("処理中..." + " " + docWithFile.Version.Name + " " + docWithFile.FileName)
		createPokemonInfoCsv(docWithFile.Doc, docWithFile.FileName, docWithFile.Index, docWithFile.SubIndex, true, docWithFile.Version)
	}
}

func createMovesCsv(version models.Version) {
	path := helpers.SavePath("csv\\move", version.Id(), "move.csv")
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"id\",\"name\",\"type\",\"power1\",\"power2\",\"pp\",\"accuracy\",\"priority\",\"damageType\",\"isDirect\",\"canProtect\"" + "\n"
		}
	}

	moves := CretaMoves(version)
	rt := reflect.New(reflect.TypeOf(models.Move{})).Elem().Type()
	for _, move := range moves {
		result += createData(move, rt)
	}

	helpers.SaveCsv(result, version.Id(), "moves.csv")
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

func createPokemonInfoCsv(page *goquery.Document, id string, index int, subIndex int, isDefault bool, version models.Version) {
	pokemon := CreatePokemon(page, id, index, subIndex, isDefault, version)
	pokemonNames := CreatePokemonName(page, pokemon, version)
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"pokemonId\",\"localLanguageId\",\"name\",\"formName\"" + "\n"
		}
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"pokemonId\",\"hp\",\"attack\",\"defense\",\"spAttack\",\"spDefense\",\"speed\"" + "\n"
		}
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"pokemonId\",\"typeId\"" + "\n"
		}
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"pokemonId\",\"abilityName\",\"isHidden\"" + "\n"
		}
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"pokemonId\",\"moveName\"" + "\n"
		}
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			result = "\"pokemonId\",\"evolutionChainId\",\"order\"" + "\n"
		}
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
	path := helpers.SavePath("csv", version, fileName)
	var result string
	if _, err := os.Stat(path); err != nil {
		if err != nil && os.IsNotExist(err) {
			result = "\"id\",\"no\",\"height\",\"weight\",\"order\",\"isDefault\"" + "\n"
		}
	}

	rt := reflect.New(reflect.TypeOf(models.Pokemon{})).Elem().Type()
	for _, row := range list {
		result += createData(row, rt)
	}

	helpers.SaveCsv(result, version, fileName)
}

func getForms(page *goquery.Document, index int) []string {
	if index == 773 {
		return []string{
			"773a", // かくとう
			"773b", // ひこう
			"773c", // どく
			"773d", // じめん
			"773e", // いわ
			"773f", // むし
			"773g", // ゴースト
			"773h", // はがね
			"773i", // ほのお
			"773j", // みず
			"773k", // くさ
			"773l", // でんき
			"773m", // エスパー
			"773n", // こおり
			"773o", // ドラゴン
			"773p", // あく
			"773q", // フェアリー
		}
	}

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
