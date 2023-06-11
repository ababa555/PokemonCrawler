package scraper_html

import (
	"PokemonCrawler/app/models"
	"PokemonCrawler/app/pkg/helpers"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractionNumberFromString() func(target string) int {
	return func(target string) int {
		if target == "" {
			return 0
		}
		r := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		numberStrings := r.FindAllString(target, -1)
		n, _ := strconv.Atoi(numberStrings[0])
		return n
	}
}

// CretaMoves 技の一覧を作成します
func CretaMoves(version models.Version) models.Moves {
	var moves models.Moves
	no := 1

	var urlParams []string
	urlParams = []string{""}
	if version.Name == "sm" {
		urlParams = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	}
	for _, urlParam := range urlParams {
		page := helpers.VisitMoveHtml(version.Id(), urlParam)
		page.Find("#contents > div:nth-child(4) > div").Each(func(index int, s *goquery.Selection) {
			var name string
			var typeAsString string
			var damageType string
			var power string
			var power2 string
			var pp int
			var accuracy string
			var priority int
			var isDirect bool
			var canProtect bool
			var preC1 bool
			s.Find("#contents > div:nth-child(4) > div > table.center > tbody > tr").Each(func(index int, s1 *goquery.Selection) {
				className, _ := s1.Attr("class")

				// sort_trはソードシールドの技リストのページ対応
				if className == "c1" || className == "sort_tr" {
					i := 2
					name = s1.Find("td.left").Text()
					typeAsString = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
					i++

					damageType = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
					i++

					power = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
					i++

					// Z技 or ダイマックス
					if version.Name == "swsh" || version.Name == "sm" {
						power2 = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
						i++
					}

					accuracy = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
					i++

					pp, _ = strconv.Atoi(s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text())
					i++

					if version.Name == "swsh" || version.Name == "sm" {
						if strings.Contains(s1.Find("td:nth-child("+strconv.Itoa(i)+")").Text(), "○") {
							isDirect = true
						} else {
							isDirect = false
						}
						i++
					}

					if strings.Contains(s1.Find("td:nth-child("+strconv.Itoa(i)+")").Text(), "○") {
						canProtect = true
					} else {
						canProtect = false
					}
					i++
					preC1 = true
					return
				} else if preC1 {
					r := regexp.MustCompile("(優先度:([+\\-0-9]+))")
					match := r.FindString(s1.Text())
					e := extractionNumberFromString()
					priority = e(match)
					preC1 = false
				}

				if name == "" {
					return
				}
				typeID, _ := models.TypeAsEnum(typeAsString)

				move := models.Move{
					ID:         version.No + "-" + strconv.Itoa(no),
					Name:       name,
					TypeID:     typeID,
					Power:      power,
					Power2:     power2,
					Pp:         pp,
					Accuracy:   accuracy,
					Priority:   priority,
					DamageType: damageType, // 1（ステータス変化）２（物理技）3（特殊技）
					IsDirect:   isDirect,
					CanProtect: canProtect,
				}
				moves = append(moves, move)
				no++
				name = ""
			})
		})
	}
	return moves
}
