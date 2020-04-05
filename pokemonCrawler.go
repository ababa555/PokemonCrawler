package main

import (
	"bufio"
	"flag"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"./pkg/models"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

var directory = `C:\temp`

func createSearchURL(version string, no string) string {
	baseURL, _ := url.Parse(strings.Join([]string{
		"https://yakkun.com/",
		version,
		"/",
		"zukan",
		"/",
		no,
	}, ""))

	return baseURL.String()
}

func visit(searchurl string) *goquery.Document {
	resp, err := http.Get(searchurl)
	if err != nil {
		return nil
	}

	// EUC_JPからUTF8に変換
	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.EUCJP.NewDecoder())
	doc, _ := goquery.NewDocumentFromReader(utfBody)

	return doc
}

func extractionWeightOrHeightFromString() func(target string) int {
	return func(target string) int {
		if target == "" {
			return 0
		}
		pattern := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		numberStrings := pattern.FindAllString(target, -1)
		numberFloat, _ := strconv.ParseFloat(numberStrings[0], 32)
		numberInt, _ := strconv.Atoi(strconv.FormatFloat(numberFloat*10, 'g', 4, 64))
		return numberInt
	}
}

func createPokemon(page *goquery.Document, index int, isDefault bool, hasGaralNo bool) models.Pokemon {
	adjust := map[bool]int{true: 1, false: 0}[hasGaralNo]
	no := page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(4+adjust) + ") > td:nth-child(2)").First().Text()
	heightAsString := page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(5+adjust) + ") > td:nth-child(2)").First().Text()
	weightAsString := page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(6+adjust) + ") > td:nth-child(2)").First().Text()
	order := index

	e := extractionWeightOrHeightFromString()
	height := e(heightAsString)
	weight := e(weightAsString)

	pokemon := models.Pokemon{
		ID:        no,
		No:        no,
		Height:    height,
		Weight:    weight,
		Order:     order,
		IsDefault: isDefault,
	}

	return pokemon
}

func createPokemonNames(page *goquery.Document, pokemon models.Pokemon, hasGaralNo bool) models.PokemonNames {
	var pokemonNames models.PokemonNames

	name := page.Find("#base_anchor > table > tbody > tr.head > th").First().Text()
	pokemonName := models.PokemonName{
		PokemonID:       pokemon.ID,
		LocalLanguageID: 1,
		Name:            name,
	}
	pokemonNames = append(pokemonNames, pokemonName)

	adjust := map[bool]int{true: 1, false: 0}[hasGaralNo]
	nameEn := page.Find("#base_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(8+adjust) + ") > td:nth-child(2) > ul > li").First().Text()
	pokemonNameEn := models.PokemonName{
		PokemonID:       pokemon.ID,
		LocalLanguageID: 2,
		Name:            nameEn,
	}
	pokemonNames = append(pokemonNames, pokemonNameEn)

	return pokemonNames
}

func createPokemonStats(page *goquery.Document, pokemon models.Pokemon) models.PokemonStats {
	hp, _ := strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(2) > td.left").Text()))
	attack, _ := strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(3) > td.left").Text()))
	defense, _ := strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(4) > td.left").Text()))
	spAttack, _ := strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(5) > td.left").Text()))
	spDefense, _ := strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(6) > td.left").Text()))
	speed, _ := strconv.Atoi(removeNbsp(page.Find("#stats_anchor > table > tbody > tr:nth-child(7) > td.left").Text()))

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

func createPokemonTypes(page *goquery.Document, pokemon models.Pokemon) models.PokemonTypes {
	var pokemonTypes models.PokemonTypes
	page.Find("#base_anchor > table > tbody > tr:nth-child(7) > td > ul").Each(func(index int, s *goquery.Selection) {
		s.Find("li > a > img").EachWithBreak(func(_ int, s2 *goquery.Selection) bool {
			typeAsString, _ := s2.Attr("alt")
			typeID, _ := models.TypeAsEnum(typeAsString)

			pokemonType := models.PokemonType{
				PokemonID: pokemon.ID,
				TypeID:    typeID,
			}
			pokemonTypes = append(pokemonTypes, pokemonType)
			return true
		})
	})
	return pokemonTypes
}

func createPokemonAbilities(page *goquery.Document, pokemon models.Pokemon) models.PokemonAbilities {
	const start = 34
	adjust := 0
	for {
		text := page.Find("#stats_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(start+adjust) + ") > th").Text()
		if !strings.Contains(text, "特性(とくせい)") {
			adjust++
			continue
		}
		break
	}
	// 特性
	var pokemonAbilities models.PokemonAbilities
	for {
		text := page.Find("#stats_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(start+1+adjust) + ") > td.c1 > a").Text()
		if text == "" {
			break
		}
		pokemonAbility := models.PokemonAbility{
			PokemonID:   pokemon.ID,
			AbilityName: text,
			IsHidden:    false,
		}
		pokemonAbilities = append(pokemonAbilities, pokemonAbility)
	}
	// 隠れ特性
	for {
		text := page.Find("#stats_anchor > table > tbody > tr:nth-child(" + strconv.Itoa(start+3+adjust) + ") > td.c1 > a").Text()
		if text == "" {
			break
		}
		pokemonAbility := models.PokemonAbility{
			PokemonID:   pokemon.ID,
			AbilityName: text,
			IsHidden:    true,
		}
		pokemonAbilities = append(pokemonAbilities, pokemonAbility)
	}
	return pokemonAbilities
}

func removeNbsp(target string) string {
	return strings.Replace(target, "\u00A0", "", 1)
}

func saveCsv(data string, filename string) {
	exe, err := os.Executable()
	current := filepath.Dir(exe)
	path := filepath.Join(current, "pkg", "files", filename)

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
			value = reflect.Indirect(r).FieldByName(f.Name).String()
		} else if f.Type.String() == "int" {
			value = strconv.FormatInt(reflect.Indirect(r).FieldByName(f.Name).Int(), 10)
		} else {
			value = reflect.Indirect(r).FieldByName(f.Name).String()
		}

		if i == 0 {
			result = value
		} else {
			result = result + "," + value
		}
	}
	return result + "\n"
}

func main() {
	flag.StringVar(&directory, "o", `C:\temp`, "画像の保存先")
	flag.Parse()

	versions := [3]string{
		"swsh", "pika_vee", "sm",
	}

	// all1 := make([]typefile.Pokemon, 0)
	// all2 := make([]typefile.PokemonName, 0)

	// for _, version := range versions {
	// 	moves := logics.CretaMoves(version)
	// }
	var moves models.Moves
	move := models.Move{}
	moves = append(moves, move)
	moves = append(moves, move)

	var result string
	rt := reflect.New(reflect.TypeOf(models.Move{})).Elem().Type()
	for _, move := range moves {
		result += createData(move, rt)
	}
	saveCsv(result, "moves.csv")

	index := 1
	for _, version := range versions {
		for {
			// ここに無限ループで実行する処理を記述
			searchNo := "n" + strconv.Itoa(index)
			searchurl := createSearchURL(version, searchNo)
			page := visit(searchurl)

			// ガラル地方のポケモンかどうかを調べる
			var hasGaralNo = false
			if version == "swsh" {
				text := page.Find("#base_anchor > table > tbody > tr:nth-child(4) > td.c1").Text()
				if text == "ガラルNo." {
					hasGaralNo = true
				}
			}

			pokemon := createPokemon(page, index, true, hasGaralNo)
			// pokemonNames := createPokemonNames(page, pokemon, hasGaralNo)
			// pokemonStats := createPokemonStats(page, pokemon)
			// pokemonTypes := createPokemonTypes(page, pokemon)
			// pokemonAbilities := createPokemonAbilities(page, pokemon)

			// pokemon_moves
			var pokemonMoves models.PokemonMoves
			page.Find("#move_anchor").Each(func(index int, s *goquery.Selection) {
				s.Find("#move_list > tbody > tr").EachWithBreak(func(index int, s1 *goquery.Selection) bool {
					name := s1.Find("td.move_name_cell > a").Text()
					title := s1.Find("#move_list > tbody > tr:nth-child(" + strconv.Itoa(index+1) + ") > th").Text()
					if strings.Contains(title, "過去作でしか覚えられない技") {
						return false
					}
					if name != "" {
						print(name)
						pokemonMove := models.PokemonMove{
							PokemonID: pokemon.ID,
							Version:   version,
							MoveName:  name,
						}
						pokemonMoves = append(pokemonMoves, pokemonMove)
					}
					return true
				})
			})

			// page.Find("#base_anchor > table > tbody > tr.head > th").Each(func(index int, s *goquery.Selection) {
			// 	title := s.Text()
			// })

			// add
			// all1 = append(all1, pokemon)
			// all2 = append(all2, pokemonNames...)

			//fmt.Println(pokemon, pokemonNames, pokemonStats, pokemonTypes, pokemonAbilities)
			index++
			return
		}
	}
}
