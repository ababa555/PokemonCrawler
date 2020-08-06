package models

import "errors"

// Pokemon 基本情報
type Pokemon struct {
	ID        string `csv:"id"` // n1、n3mなど
	No        string `csv:"no"` // ぜんこくNo.
	Height    int    `csv:"height"`
	Weight    int    `csv:"weight"`
	Order     int    `csv:"order"`
	IsDefault bool   `csv:"isDefault"`
}

// Pokemons 基本情報のリスト
type Pokemons []Pokemon

// PokemonName ポケモン名
type PokemonName struct {
	PokemonID       string `csv:"pokemonId"`
	LocalLanguageID int    `csv:"localLanguageId"`
	Name            string `csv:"name"`
	FormName        string `csv:"formName"`
}

// PokemonNames ポケモン名のリスト
type PokemonNames []PokemonName

// PokemonEvolutionChain ポケモンの進化情報
type PokemonEvolutionChain struct {
	PokemonID        string `csv:"pokemonId"`
	EvolutionChainID string `csv:"evolutionChainId"`
	Order            int    `csv:"order"`
}

// PokemonEvolutionChains ポケモンの進化情報のリスト
type PokemonEvolutionChains []PokemonEvolutionChain

// PokemonStats ポケモンのステータス
type PokemonStats struct {
	PokemonID string `csv:"pokemonId"`
	Hp        int    `csv:"hp"`
	Attack    int    `csv:"attack"`
	Defense   int    `csv:"defense"`
	SpAttack  int    `csv:"spAttack"`
	SpDefense int    `csv:"spDefense"`
	Speed     int    `csv:"speed"`
}

// PokemonStatses 基本情報のリスト
type PokemonStatses []PokemonStats

// PokemonType ポケモンのタイプ
type PokemonType struct {
	PokemonID string `csv:"pokemonId"`
	TypeID    Type   `csv:"typeId"`
}

// PokemonTypes ポケモンのタイプのリスト
type PokemonTypes []PokemonType

// PokemonAbility ポケモンの特性
type PokemonAbility struct {
	PokemonID   string `csv:"pokemonId"`
	AbilityName string `csv:"abilityName"`
	IsHidden    bool   `csv:"isHidden"`
}

// PokemonAbilities ポケモンの特性のリスト
type PokemonAbilities []PokemonAbility

// PokemonMove ポケモンの技
type PokemonMove struct {
	PokemonID string `csv:"pokemonId"`
	MoveName  string `csv:"moveName"`
}

// PokemonMoves ポケモンの技のリスト
type PokemonMoves []PokemonMove

// Move 技
type Move struct {
	ID         int // index
	Name       string
	TypeID     Type
	Power      string
	Power2     string // Z技 or ダイマックス
	Pp         int
	Accuracy   string // 命中率
	Priority   int    // 優先度
	DamageType string // 1（ステータス変化）２（物理技）3（特殊技）
	IsDirect   bool   // 直接攻撃か？
	CanProtect bool   // 守るできるか？
}

// Moves 技のリスト
type Moves []Move

// Type ポケモンのタイプ
type Type int

const (
	// Normal ノーマル
	Normal = iota
	// Fighting かくとう
	Fighting
	// Flying ひこう
	Flying
	// Poison どく
	Poison
	// Ground じめん
	Ground
	// Rock いわ
	Rock
	// Bug むし
	Bug
	// Ghost ゴースト
	Ghost
	// Steel はがね
	Steel
	// Fire ほのお
	Fire
	// Water みず
	Water
	// Grass くさ
	Grass
	// Electric でんき
	Electric
	// Psychic エスパー
	Psychic
	// Ice こおり
	Ice
	// Dragon ドラゴン
	Dragon
	// Dark あく
	Dark
	// Fairy フェアリー
	Fairy
)

// TypeAsEnum 文字列をType(ポケモンのタイプ)に変換します。
func TypeAsEnum(typeAsString string) (Type, error) {
	switch typeAsString {
	case "ノーマル":
		return Normal, nil
	case "かくとう":
		return Fighting, nil
	case "ひこう":
		return Flying, nil
	case "どく":
		return Poison, nil
	case "じめん":
		return Ground, nil
	case "いわ":
		return Rock, nil
	case "むし":
		return Bug, nil
	case "ゴースト":
		return Ghost, nil
	case "はがね":
		return Steel, nil
	case "ほのお":
		return Fire, nil
	case "みず":
		return Water, nil
	case "くさ":
		return Grass, nil
	case "でんき":
		return Electric, nil
	case "エスパー":
		return Psychic, nil
	case "こおり":
		return Ice, nil
	case "ドラゴン":
		return Dragon, nil
	case "あく":
		return Dark, nil
	case "フェアリー":
		return Fairy, nil
	default:
		return 0, errors.New("invalid type")
	}
}
