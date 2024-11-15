package piece

import "fmt"

type Pattern uint8
type Face uint8

const (
	Blue     = 1
	BlueStar = 1

	Border = 2

	Purple     = 3
	PinkCircle = 3

	Yellow      = 4
	PinkTrident = 4

	White        = 5
	YellowCircle = 5

	Pink              = 6
	RedYellowTriangle = 6

	Green               = 7
	YellowGreenTriangle = 7

	Red        = 8
	RedTrident = 8
)

const (
	North = 0
	East  = 1
	South = 2
	West  = 3
)

func (p Pattern) String() string {
	switch p {
	case 0:
		return "--"
	case Blue:
		return "Blue"
	case Border:
		return "Border"
	case Purple:
		return "Purple"
	case Yellow:
		return "Yellow"
	case White:
		return "White"
	case Pink:
		return "Pink"
	case Green:
		return "Green"
	case Red:
		return "Red"
	default:
		return "Unknown"
	}
}

func (f Face) String() string {
	switch f {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return "Unknown"
	}
}

type Piece struct {
	number int
	placed bool
	sticky bool
	north  Pattern
	east   Pattern
	south  Pattern
	west   Pattern
}

type PiecePlacement struct {
	piece       *Piece
	orientation Face
	north       Pattern
	east        Pattern
	south       Pattern
	west        Pattern
}

var nextPieceNumber int = 1
var firstCorner bool = true

func New(north Pattern, east Pattern, south Pattern, west Pattern) Piece {
	piece := Piece{north: north, east: east, south: south, west: west, number: nextPieceNumber, placed: false}
	borders := 0
	if north == Border {
		borders += 1
	}
	if east == Border {
		borders += 1
	}
	if south == Border {
		borders += 1
	}
	if west == Border {
		borders += 1
	}

	if firstCorner && borders == 2 {
		piece.sticky = true
		firstCorner = false
	}

	nextPieceNumber += 1
	return piece
}

func (p Piece) String() string {
	return fmt.Sprintf("Piece %d: N=%s E=%s S=%s W=%s", p.number, p.north, p.east, p.south, p.west)
}

func (p Piece) Rotations() [4]PiecePlacement {
	return [4]PiecePlacement{
		{&p, North, p.north, p.east, p.south, p.west},
		{&p, East, p.east, p.south, p.west, p.north},
		{&p, South, p.south, p.west, p.north, p.east},
		{&p, West, p.west, p.north, p.east, p.south},
	}
}

func (pp PiecePlacement) Keys() [11]int {

	var N = int(pp.north)
	var E = int(pp.east) * 9
	var S = int(pp.south) * 81
	var W = int(pp.west) * 729

	return [11]int{
		N + E + S + W,
		E + S + W,
		N + S + W,
		N + E + W,
		N + E + S,
		N + E,
		N + S,
		N + W,
		E + S,
		E + W,
		S + W,
	}
}
