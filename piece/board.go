package piece

import (
	"fmt"
)

const BOARD_SIZE = 6
const DEBUG = 0
const INFO = 1
const WARN = 2
const ERROR = 3
const LOG_LEVEL = 2

var MovementCount = 0

type Tile struct {
	backtracking_queue []*PiecePlacement
	placed_piece       *PiecePlacement
	north_restriction  Pattern
	east_restriction   Pattern
	south_restriction  Pattern
	west_restriction   Pattern
}

func (t *Tile) LookupKey() int {
	return int(t.north_restriction) | int(t.east_restriction)<<8 | int(t.south_restriction)<<16 | int(t.west_restriction)<<24
}

func (t *Tile) Dequeue() *PiecePlacement {
	result := t.backtracking_queue[0]
	t.backtracking_queue = t.backtracking_queue[1:]
	return result
}

type Coordinate struct {
	x int
	y int
}

var IterationOrder = [36]Coordinate{
	{x: 1, y: 1},
}

type Board struct {
	tiles        [36]*Tile
	currentPiece int
}

func (b Board) String() string {
	str := "<!doctype html>\n<html>\n<head>  <link rel=\"stylesheet\" type=\"text/css\" href=\"style.css\" >\n</head>\n<body>\n  <table class=\"board\">\n"
	idx := 0
	for y := 1; y <= BOARD_SIZE; y++ {
		str += "    <tr>\n"
		for x := 1; x <= BOARD_SIZE; x++ {
			pp := b.tiles[idx].placed_piece
			str += fmt.Sprintf("      <td><img src=\"piece%d.png\" class=\"img img-%s\" /></td>\n", pp.piece.number, pp.orientation.String()[:1])
			idx += 1
		}
		str += "    </tr>\n"
	}
	return str + "  </table>\n</body>\n</html>"
}

func NewBoard(pieceLookup map[int]*PiecePlacementLookup) Board {
	board := Board{}
	board.currentPiece = 0

	// set up tiles and border restrictions
	for x := 1; x <= BOARD_SIZE; x++ {
		for y := 1; y <= BOARD_SIZE; y++ {
			idx := Coordinate{x: x, y: y}.AsIndex()
			board.tiles[idx] = &Tile{}
			if x == 1 {
				board.tiles[idx].west_restriction = Border
			}
			if x == BOARD_SIZE {
				board.tiles[idx].east_restriction = Border
			}
			if y == 1 {
				board.tiles[idx].north_restriction = Border
			}
			if y == BOARD_SIZE {
				board.tiles[idx].south_restriction = Border
			}
		}
	}

	board.currentPiece += 1
	board.Place(pieceLookup[board.tiles[0].LookupKey()].GetPieces()[0], Coordinate{x: 1, y: 1}, pieceLookup)

	return board
}

func (b *Board) IsSolved() bool {
	return b.currentPiece == BOARD_SIZE*BOARD_SIZE
}

func BuildLookup(pieces [36]Piece) map[int]*PiecePlacementLookup {
	pieceLookup := make(map[int]*PiecePlacementLookup)
	for _, p := range pieces {
		for _, pp := range p.Rotations() {
			for _, k := range pp.Keys() {
				if _, ok := pieceLookup[k]; !ok {
					pieceLookup[k] = NewPiecePlacementLookup()
				}
				sizeBefore := len(pp.piece.Lookups)
				pp.piece.Lookups[k] = pieceLookup[k]
				sizeAfter := len(pp.piece.Lookups)
				if sizeAfter > sizeBefore {
					pp.piece.LookupsList = append(pp.piece.LookupsList, pieceLookup[k])
				}
				pieceLookup[k].Add(&pp)
			}
		}
	}
	return pieceLookup
}

func (b *Board) GetNextCoordinate(pieceLookup map[int]*PiecePlacementLookup) Coordinate {
	minChoices := 10000
	minTile := -1

	for i, t := range b.tiles {
		if t.placed_piece == nil {
			restrictions := 0
			if t.east_restriction != 0 {
				restrictions += 1
			}
			if t.west_restriction != 0 {
				restrictions += 1
			}
			if t.north_restriction != 0 {
				restrictions += 1
			}
			if t.south_restriction != 0 {
				restrictions += 1
			}

			if restrictions < 2 {
				continue
			}

			lookupKey := t.LookupKey()
			lookup := pieceLookup[lookupKey]
			choices := 0
			if lookup != nil {
				choices = lookup.count
			}
			if choices < minChoices {
				minChoices = choices
				minTile = i
			}
		}
	}

	IterationOrder[b.currentPiece] = Coordinate{x: minTile%6 + 1, y: minTile/6 + 1}
	return IterationOrder[b.currentPiece]
}

func (b *Board) PlaceNext(pieceLookup map[int]*PiecePlacementLookup) bool {
	// get coordinate of next tile
	coord := b.GetNextCoordinate(pieceLookup)
	idx := coord.AsIndex()
	tile := b.tiles[idx]
	backtracking_queue := make([]*PiecePlacement, 0, 1)

	// get matching pieces
	matchingPieces := pieceLookup[tile.LookupKey()]
	if LOG_LEVEL <= INFO {
		ordIndicator := "th"
		if (b.currentPiece+1)%10 == 1 {
			ordIndicator = "st"
		} else if (b.currentPiece+1)%10 == 2 {
			ordIndicator = "nd"
		} else if (b.currentPiece+1)%10 == 3 {
			ordIndicator = "rd"
		}
		fmt.Printf("Will attempt to place %d%s piece at %d, %d\n", b.currentPiece+1, ordIndicator, coord.x, coord.y)
	}

	if LOG_LEVEL <= DEBUG {
		fmt.Printf("  Restrictions for this tile: N=%s E=%s S=%s W=%s\n", tile.north_restriction, tile.east_restriction, tile.south_restriction, tile.west_restriction)
		fmt.Printf("  Matching piece placements: %d\n", len(backtracking_queue))
	}

	if matchingPieces == nil || matchingPieces.count == 0 {
		if LOG_LEVEL <= INFO {
			fmt.Printf("  None of the remaining pieces is a match, backtracking\n")
		}
		return false
	}
	b.currentPiece += 1
	tile.backtracking_queue = matchingPieces.GetPieces()

	nextCandidate := tile.Dequeue()
	if LOG_LEVEL <= INFO {
		fmt.Printf("  Placing piece %d facing %s (%s, %s, %s, %s)\n", nextCandidate.piece.number, nextCandidate.orientation, nextCandidate.north, nextCandidate.east, nextCandidate.south, nextCandidate.west)
	}
	b.Place(nextCandidate, coord, pieceLookup)
	return true
}

func (b *Board) Place(pp *PiecePlacement, xy Coordinate, pieceLookup map[int]*PiecePlacementLookup) {
	MovementCount += 1
	idx := xy.AsIndex()
	tile := b.tiles[idx]
	tile.placed_piece = pp
	pp.MarkUsed()

	// update restrictions on neighboring tiles and check if it doesn't create an unsolvable situation
	if xy.x > 1 {
		left_tile := b.tiles[idx-1]
		left_tile.east_restriction = pp.west
	}
	if xy.x < BOARD_SIZE {
		right_tile := b.tiles[idx+1]
		right_tile.west_restriction = pp.east
	}
	if xy.y > 1 {
		top_tile := b.tiles[idx-BOARD_SIZE]
		top_tile.south_restriction = pp.north
	}
	if xy.y < BOARD_SIZE {
		bottom_tile := b.tiles[idx+BOARD_SIZE]
		bottom_tile.north_restriction = pp.south
	}
}

func (b *Board) Unplace(xy Coordinate) {
	MovementCount += 1
	idx := xy.AsIndex()
	tile := b.tiles[idx]
	tile.placed_piece.MarkUnused()
	tile.placed_piece = nil

	// update restrictions on neighboring tiles
	if xy.x > 1 {
		left_tile := b.tiles[idx-1]
		left_tile.east_restriction = 0
	}
	if xy.x < BOARD_SIZE {
		right_tile := b.tiles[idx+1]
		right_tile.west_restriction = 0
	}
	if xy.y > 1 {
		top_tile := b.tiles[idx-BOARD_SIZE]
		top_tile.south_restriction = 0
	}
	if xy.y < BOARD_SIZE {
		bottom_tile := b.tiles[idx+BOARD_SIZE]
		bottom_tile.north_restriction = 0
	}
}

func (b *Board) Backtrack(pieceLookup map[int]*PiecePlacementLookup) bool {
	// start by removing the last placed piece
	currentCoordinate := IterationOrder[b.currentPiece-1]
	if LOG_LEVEL <= INFO {
		fmt.Printf("Removing current placed piece at %d, %d\n", currentCoordinate.x, currentCoordinate.y)
	}
	b.Unplace(currentCoordinate)

	// check the next piece in the backtracking queue
	tile := b.tiles[currentCoordinate.AsIndex()]

	if len(tile.backtracking_queue) == 0 {
		if LOG_LEVEL <= DEBUG {
			fmt.Printf("  Backtracking further as there are no more candidates for %d, %d\n", currentCoordinate.x, currentCoordinate.y)
		}
		tile.backtracking_queue = nil
		b.currentPiece -= 1
		if b.currentPiece == 0 {
			return false
		} else {
			return b.Backtrack(pieceLookup)
		}
	}

	nextCandidate := tile.Dequeue()
	if LOG_LEVEL <= INFO {
		fmt.Printf("  Placing piece %d facing %s\n", nextCandidate.piece.number, nextCandidate.orientation)
	}
	b.Place(nextCandidate, currentCoordinate, pieceLookup)
	return true
}

func (xy Coordinate) AsIndex() int {
	return (xy.y-1)*6 + xy.x - 1
}
