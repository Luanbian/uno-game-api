// Package game card provides a card logic and control
package game

type (
	Color string
	Type  string
)

const (
	Red    Color = "red"
	Green  Color = "green"
	Blue   Color = "blue"
	Yellow Color = "yellow"
	None   Color = "none"

	Number   Type = "number"
	Jump     Type = "jump"
	Inverter Type = "inverter"
	Plustwo  Type = "plusTwo"
	Plusfour Type = "plusFour"
	Joker    Type = "joker"

	ColorSelect Type = "colorSelect"
)

type Card struct {
	Color Color
	Type  Type
	Value int // 0-9 para Number, -1 para especiais
}
