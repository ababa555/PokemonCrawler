package main

import (
	"PokemonCrawler/app/crawler_image"
	"PokemonCrawler/app/crawler_info"
)

func main() {

	println("処理を開始します...")

	// ポケモン徹底攻略をクローラ
	crawler_info.Run()

	// ポケモン図鑑(https://zukan.pokemon.co.jp)から画像をクローラ
	crawler_image.Run()
}
