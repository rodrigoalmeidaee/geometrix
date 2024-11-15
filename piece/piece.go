package piece

import "fmt"

type Pattern uint8
type Face uint8

const (
	Blue     = 1
	BlueStar = 1

	Border = 2

	Purple     = 4
	PinkCircle = 4

	Yellow      = 8
	PinkTrident = 8

	White        = 16
	YellowCircle = 16

	Pink              = 32
	RedYellowTriangle = 32

	Green               = 64
	YellowGreenTriangle = 64

	Red        = 128
	RedTrident = 128
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

func New(north Pattern, east Pattern, south Pattern, west Pattern) Piece {
	piece := Piece{north: north, east: east, south: south, west: west, number: nextPieceNumber, placed: false}
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

func (pp PiecePlacement) Keys() [16]int {

	var N = int(pp.north)
	var E = int(pp.east) << 8
	var S = int(pp.south) << 16
	var W = int(pp.west) << 24

	return [16]int{
		N | E | S | W,
		E | S | W,
		N | S | W,
		N | E | W,
		N | E | S,
		N | E,
		N | S,
		N | W,
		E | S,
		E | W,
		S | W,
		N,
		E,
		S,
		W,
		0,
	}
}
