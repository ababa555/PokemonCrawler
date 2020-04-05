package models

import "errors"

// Pokemon 基本情報
type Pokemon struct {
	ID        string // n1、n3mなど
	No        string // ぜんこくNo.
	Height    int
	Weight    int
	Order     int
	IsDefault bool
}
type pokemons []Pokemon

// PokemonName ポケモン名
type PokemonName struct {
	PokemonID       string
	LocalLanguageID int
	Name            string
}

// PokemonNames ポケモン名のリスト
type PokemonNames []PokemonName

// PokemonEvolutionChain is struct
type PokemonEvolutionChain struct {
	PokemonID        string
	EvolutionChainID string
	Name             string
}

// PokemonStats ポケモンのステータス
type PokemonStats struct {
	PokemonID string
	Hp        int
	Attack    int
	Defense   int
	SpAttack  int
	SpDefense int
	Speed     int
}

// PokemonType ポケモンのタイプ
type PokemonType struct {
	PokemonID string
	TypeID    Type
}

// PokemonTypes ポケモンのタイプのリスト
type PokemonTypes []PokemonType

// PokemonAbility ポケモンの特性
type PokemonAbility struct {
	PokemonID   string
	AbilityName string
	IsHidden    bool
}

// PokemonAbilities ポケモンの特性のリスト
type PokemonAbilities []PokemonAbility

// PokemonMove ポケモンの技
type PokemonMove struct {
	PokemonID string
	Version   string
	MoveName  string
}

// PokemonMoves ポケモンの技のリスト
type PokemonMoves []PokemonMove

// Move 技
type Move struct {
	ID         int // index
	Version    string
	Name       string
	TypeID     Type
	Power      string
	Power2     string // Z技 or ダイマックス
	Pp         int
	Accuracy   string
	Priority   int
	DamageType string // 1（ステータス変化）２（物理技）3（特殊技）
	IsDirect   bool
	CanProtect bool
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
		return 0, errors.New("division by zero")
	}
}
