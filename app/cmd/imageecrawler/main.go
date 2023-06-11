package main

import "PokemonCrawler/app/pkg/crawler/crawler_image"

func main() {

	println("処理を開始します...")

	// ポケモン図鑑(https://zukan.pokemon.co.jp)から画像をクローラ
	crawler_image.Run()
}
