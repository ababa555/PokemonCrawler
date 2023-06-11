package crawler_html

import (
	"PokemonCrawler/app/models"
	"PokemonCrawler/app/pkg/downloader"
	"PokemonCrawler/app/pkg/helpers"
	"log"
	"strconv"
)

func Run() {
	versions := []models.Version{
		*models.NewVersion("1", "sm"),
		*models.NewVersion("2", "pika_vee"),
		*models.NewVersion("3", "swsh"),
		*models.NewVersion("4", "sv"),
	}

	for _, version := range versions {
		downloadMoveHtml(version)
	}

	for _, version := range versions {
		downloadPokemonHtml(version)
	}
}

func downloadMoveHtml(version models.Version) {
	var urlParams []string
	urlParams = []string{""}
	if version.Name == "sm" {
		urlParams = []string{"?c=0", "?c=1", "?c=2", "?c=3", "?c=4", "?c=5", "?c=6", "?c=7", "?c=8", "?c=9"}
	}

	var urls []string
	for _, urlParam := range urlParams {
		searchurl := helpers.CreateSearchURLMovePage(version.Name, urlParam)
		urls = append(urls, searchurl)
	}

	ps, err := downloader.NewWebPageScraper()
	if err != nil {
		log.Fatalf("could not create page scraper: %v", err)
	}

	if err := ps.ScrapeMoveHtml(urls, version.Id()); err != nil {
		log.Fatalf("could not scrape pages: %v", err)
	}
}

func downloadPokemonHtml(version models.Version) {
	index := 1
	for {
		// ここに無限ループで実行する処理を記述
		searchNo := "n" + strconv.Itoa(index)
		println("処理中... " + searchNo + "(" + version.Name + ")")

		// ダウンロード済の場合、次のポケモンへ
		// if isDownloaded(version.Id(), searchNo) {
		// 	index++
		// 	continue
		// }

		if searchNo != "n190" && searchNo != "n193" && searchNo != "n201" && searchNo != "n207" &&
			searchNo != "n265" && searchNo != "n266" && searchNo != "n267" && searchNo != "n268" &&
			searchNo != "n269" && searchNo != "n299" && searchNo != "n358" && searchNo != "n387" &&
			searchNo != "n388" && searchNo != "n389" && searchNo != "n390" && searchNo != "n391" &&
			searchNo != "n392" && searchNo != "n393" && searchNo != "n394" && searchNo != "n395" &&
			searchNo != "n399" && searchNo != "n400" && searchNo != "n408" && searchNo != "n409" &&
			searchNo != "n410" && searchNo != "n411" && searchNo != "n412" && searchNo != "n413" &&
			searchNo != "n414" && searchNo != "n424" && searchNo != "n431" && searchNo != "n432" &&
			searchNo != "n433" && searchNo != "n441" && searchNo != "n455" && searchNo != "n46" &&
			searchNo != "n469" && searchNo != "n47" && searchNo != "n472" && searchNo != "n476" &&
			searchNo != "n489" && searchNo != "n490" && searchNo != "n491" && searchNo != "n492" &&
			searchNo != "n74" && searchNo != "n75" && searchNo != "n76" {
			index++
			continue
		}

		var urls []string
		searchurl := helpers.CreateSearchURL(version.Name, searchNo)
		urls = append(urls, searchurl)

		// playwrightでURLのWebページにアクセスし、HTMLをローカルに保存
		ps, err := downloader.NewWebPageScraper()
		if err != nil {
			log.Fatalf("could not create page scraper: %v", err)
		}

		if err := ps.ScrapePokemonHtml(urls, version.Id(), version.Name, searchNo, index); err != nil {
			log.Fatalf("could not scrape pages: %v", err)
			break
		}

		index++
	}
}
