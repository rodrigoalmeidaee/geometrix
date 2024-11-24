package piece

import (
	"fmt"
	"slices"
)

type Pattern uint8
type Face uint8

const (
	North  = 0
	East   = 1
	South  = 2
	West   = 3
	Border = 8
)

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
	number      int
	placed      bool
	north       Pattern
	east        Pattern
	south       Pattern
	west        Pattern
	Lookups     map[int]*PiecePlacementLookup
	LookupsList []*PiecePlacementLookup
	sticky      bool
}

type PiecePlacement struct {
	piece       *Piece
	orientation Face
	north       Pattern
	east        Pattern
	south       Pattern
	west        Pattern
	fingerprint int
}

var nextPieceNumber int = 1
var cornersFound int = 0

func ResetPieceNumber() {
	nextPieceNumber = 1
	cornersFound = 0
}

func New(north Pattern, east Pattern, south Pattern, west Pattern) Piece {
	piece := Piece{north: north, east: east, south: south, west: west, number: nextPieceNumber, placed: false, Lookups: make(map[int]*PiecePlacementLookup), LookupsList: make([]*PiecePlacementLookup, 0, 16)}
	nextPieceNumber += 1
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
	if borders == 2 {
		if cornersFound == 0 {
			piece.sticky = true
		}
		cornersFound += 1
	}
	return piece
}

func (p Piece) String() string {
	return fmt.Sprintf("Piece %d: N=%d E=%d S=%d W=%d", p.number, p.north, p.east, p.south, p.west)
}

func (p Piece) Rotations() [4]PiecePlacement {
	placements := [4]PiecePlacement{
		{&p, North, p.north, p.east, p.south, p.west, 0},
		{&p, East, p.east, p.south, p.west, p.north, 0},
		{&p, South, p.south, p.west, p.north, p.east, 0},
		{&p, West, p.west, p.north, p.east, p.south, 0},
	}

	for i := 0; i < 4; i++ {
		placements[i].fingerprint = placements[i].Keys()[0]
	}

	return placements
}

func (pp PiecePlacement) Keys() [11]int {

	var N = int(pp.north)
	var E = int(pp.east) * M1
	var S = int(pp.south) * M2
	var W = int(pp.west) * M3

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

func GetPieceStats(pieces []Piece) map[int]int {
	pieceStats := make(map[int]int)

	// Statistics for single faces
	for _, p := range pieces {
		pieceStats[int(p.north)] += 1
		pieceStats[int(p.east)] += 1
		pieceStats[int(p.south)] += 1
		pieceStats[int(p.west)] += 1
	}

	// Statistics for corners
	for _, p := range pieces {
		pieceStats[int(p.north)+M1*int(p.east)] += 1
		pieceStats[int(p.east)+M1*int(p.south)] += 1
		pieceStats[int(p.south)+M1*int(p.west)] += 1
		pieceStats[int(p.west)+M1*int(p.north)] += 1
	}

	// Statistics for U shapes
	for _, p := range pieces {
		pieceStats[int(p.north)+M1*int(p.east)+M2*int(p.south)] += 1
		pieceStats[int(p.east)+M1*int(p.south)+M2*int(p.west)] += 1
		pieceStats[int(p.south)+M1*int(p.west)+M2*int(p.north)] += 1
		pieceStats[int(p.west)+M1*int(p.north)+M2*int(p.east)] += 1
	}

	return pieceStats
}

func GetPieceRarity(pieceStats map[int]int, piece Piece, facet string, aggregateMode string) int {
	var rarities []int = make([]int, 4)

	if facet == "face" {
		// Statistics for single faces
		rarities[0] = pieceStats[int(piece.north)]
		rarities[1] = pieceStats[int(piece.east)]
		rarities[2] = pieceStats[int(piece.south)]
		rarities[3] = pieceStats[int(piece.west)]
	} else if facet == "corner" {
		rarities[0] = pieceStats[int(piece.north)+M1*int(piece.east)]
		rarities[1] = pieceStats[int(piece.east)+M1*int(piece.south)]
		rarities[2] = pieceStats[int(piece.south)+M1*int(piece.west)]
		rarities[3] = pieceStats[int(piece.west)+M1*int(piece.north)]
	} else if facet == "u" {
		rarities[0] = pieceStats[int(piece.north)+M1*int(piece.east)+M2*int(piece.south)]
		rarities[1] = pieceStats[int(piece.east)+M1*int(piece.south)+M2*int(piece.west)]
		rarities[2] = pieceStats[int(piece.south)+M1*int(piece.west)+M2*int(piece.north)]
		rarities[3] = pieceStats[int(piece.west)+M1*int(piece.north)+M2*int(piece.east)]
	} else {
		panic(facet)
	}

	if aggregateMode == "min" {
		return slices.Min(rarities)
	} else if aggregateMode == "avg" {
		return (rarities[0] + rarities[1] + rarities[2] + rarities[3]) / 4
	} else {
		panic(aggregateMode)
	}
}

type PiecePlacementLookup struct {
	pieces           []*PiecePlacement
	count            uint8
	pieceRepetitions [401]uint8
}

func NewPiecePlacementLookup() *PiecePlacementLookup {
	return &PiecePlacementLookup{pieces: make([]*PiecePlacement, 0, 1), count: 0}
}

func (ppl *PiecePlacementLookup) Add(pp *PiecePlacement) {
	ppl.pieces = append(ppl.pieces, pp)
	ppl.pieceRepetitions[pp.piece.number] += 1
	ppl.count += 1
}

func (pp *PiecePlacement) MarkUsed() {
	pp.piece.placed = true
	for _, lookup := range pp.piece.LookupsList {
		lookup.MarkUsed(pp)
	}
}

func (ppl *PiecePlacementLookup) MarkUsed(pp *PiecePlacement) {
	ppl.count -= ppl.pieceRepetitions[pp.piece.number]
}

func (pp *PiecePlacement) MarkUnused() {
	pp.piece.placed = false
	for _, lookup := range pp.piece.LookupsList {
		lookup.MarkUnused(pp)
	}
}

func (ppl *PiecePlacementLookup) MarkUnused(pp *PiecePlacement) {
	ppl.count += ppl.pieceRepetitions[pp.piece.number]
}

func (ppl *PiecePlacementLookup) GetPieces() []*PiecePlacement {
	pieces := make([]*PiecePlacement, ppl.count)
	piecesSeen := make(map[int]bool, ppl.count)
	appended := 0
	for _, pp := range ppl.pieces {
		if !pp.piece.placed && !piecesSeen[pp.fingerprint] {
			pieces[appended] = pp
			piecesSeen[pp.fingerprint] = true
			appended += 1
		}
	}
	return pieces[:appended]
}
