package crawler_image

import (
	"PokemonCrawler/helpers"
	"regexp"
	"sort"
	"strings"
	"time"
)

func Run() {
	searchNo := "1"
	for {
		next, subList := crawler(searchNo)
		if next == "" {
			break
		}
		searchNo = next

		if len(subList) > 0 {
			var targetSubList []string
			sort.SliceStable(subList, func(i, j int) bool { return subList[i] < subList[j] })
			for _, v := range subList {
				target := strings.Replace(strings.TrimPrefix(strings.Replace(v, `"sub":`, ``, -1), `"`), `,`, "", -1)
				if target != "0" {
					targetSubList = append(targetSubList, target)
				}
			}

			for _, v := range helpers.SliceUnique(targetSubList) {
				next, _ := crawler(searchNo + "-" + v)
				if searchNo == next {
					searchNo = next
				} else {
					searchNo = next
					break
				}
			}
		} else {
			next, _ := crawler(searchNo)
			searchNo = next
		}
	}
}

func crawler(searchNo string) (string, []string) {
	page, statusCode := helpers.VisitZukan(searchNo)
	if statusCode == 200 {
		r1 := regexp.MustCompile(`"no":[^,]*",`)
		r2 := regexp.MustCompile(`"sub":[^,]*,`)
		r3 := regexp.MustCompile(`"image_m":[^,]*",`)
		r4 := regexp.MustCompile(`"next":{"no":[^,]*",`)

		noList := r1.FindAllString(page.Find("script").Text(), -1)
		subList := r2.FindAllString(page.Find("script").Text(), -1)
		imageList := r3.FindAllString(page.Find("script").Text(), -1)
		nextList := r4.FindAllString(page.Find("script").Text(), -1)

		no := strings.Replace(strings.TrimPrefix(strings.Replace(noList[0], `"no":`, ``, -1), `"`), `",`, "", -1)
		sub := strings.Replace(strings.TrimPrefix(strings.Replace(subList[0], `"sub":`, ``, -1), `"`), `,`, "", -1)
		imageURL := strings.Replace(strings.Replace(strings.TrimPrefix(strings.Replace(imageList[0], `"image_m":`, ``, -1), `"`), `",`, "", -1), `\`, "", -1)
		next := ""
		if len(nextList) > 0 {
			next = strings.Replace(strings.TrimPrefix(strings.Replace(nextList[0], `"next":{"no":`, ``, -1), `"`), `",`, "", -1)
		}

		// 画像を保存
		downloadURL := strings.TrimSuffix(imageURL, "\n")
		helpers.SaveImage(downloadURL, no+"_"+sub+".png")
		time.Sleep(2000 * time.Millisecond)

		// 同じポケモンのフォルム違い
		if searchNo == next {
			return next, subList
		}

		// 別のポケモン
		return next, make([]string, 0)
	}
	return "", nil
}
