package logics

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"../../pkg/models"

	"github.com/PuerkitoBio/goquery"
)

func createSearchURL(version string) string {
	baseURL, _ := url.Parse(strings.Join([]string{
		"https://yakkun.com/",
		version,
		"/",
		"move_list.htm",
	}, ""))

	return baseURL.String()
}

func extractionNumberFromString() func(target string) int {
	return func(target string) int {
		if target == "" {
			return 0
		}
		pattern := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		numberStrings := pattern.FindAllString(target, -1)
		n, _ := strconv.Atoi(numberStrings[0])
		return n
	}
}

// CretaMoves is
func CretaMoves(version string) models.Moves {
	var moves models.Moves

	searchurl := createSearchURL(version)
	page := visit(searchurl)
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
			if className == "c1" {
				i := 2
				name = s1.Find("td.left").Text()
				typeAsString = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
				i++

				damageType = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
				i++

				power = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
				i++

				if version == "swsh" || version == "sm" {
					power2 = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
					i++
				}

				accuracy = s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text()
				i++

				pp, _ = strconv.Atoi(s1.Find("td:nth-child(" + strconv.Itoa(i) + ")").Text())
				i++

				if version == "swsh" || version == "sm" {
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
				pattern := regexp.MustCompile("(優先度:([+\\-0-9]+))")
				match := pattern.FindString(s1.Text())
				e := extractionNumberFromString()
				priority = e(match)
				preC1 = false
			}
			if name == "" {
				return
			}
			typeID, _ := models.TypeAsEnum(typeAsString)
			move := models.Move{
				ID:         index,
				Version:    version,
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
		})
	})
	return moves
}
