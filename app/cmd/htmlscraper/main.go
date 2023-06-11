package main

import "PokemonCrawler/app/pkg/scraper/scraper_html"

func main() {

	println("処理を開始します...")

	// ポケモン徹底攻略をクローラしたhtmlを読み込んでcsvに変換
	scraper_html.Run()

	println("処理を終了しました...")
}
