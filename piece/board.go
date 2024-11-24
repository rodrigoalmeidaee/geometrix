package piece

import (
	"fmt"
	"time"
)

const BOARD_SIZE = 20
const DEBUG = 0
const INFO = 1
const WARN = 2
const ERROR = 3
const LOG_LEVEL = 2

const PATTERN_COUNT = 34
const M1 = (PATTERN_COUNT + 1)
const M2 = (M1 * (PATTERN_COUNT + 1))
const M3 = (M2 * (PATTERN_COUNT + 1))
const LOOKUP_SIZE = 1500625

var MovementCount = int64(0)

type Tile struct {
	backtracking_queue []*PiecePlacement
	placed_piece       *PiecePlacement
	north_restriction  Pattern
	east_restriction   Pattern
	south_restriction  Pattern
	west_restriction   Pattern
	interest_queue_pos int
	restriction_count  int
	xy                 Coordinate
}

func (t *Tile) LookupKey() int {
	return int(t.north_restriction) + int(t.east_restriction)*M1 + int(t.south_restriction)*M2 + int(t.west_restriction)*M3
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

var IterationOrder = [BOARD_SIZE * BOARD_SIZE]Coordinate{
	{x: 1, y: 1},
}

type Board struct {
	tiles              [BOARD_SIZE * BOARD_SIZE]*Tile
	currentPiece       int
	pieceLookup        [LOOKUP_SIZE]*PiecePlacementLookup
	maxPlacedPieces    int
	solveStart         int64
	interest_queue     [BOARD_SIZE * BOARD_SIZE]*Tile
	interest_queue_len int
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

func NewBoard(pieces []Piece) Board {
	board := Board{}
	board.currentPiece = 0
	board.pieceLookup = BuildLookup(pieces)

	// set up tiles and border restrictions
	for x := 1; x <= BOARD_SIZE; x++ {
		for y := 1; y <= BOARD_SIZE; y++ {
			xy := Coordinate{x: x, y: y}
			idx := xy.AsIndex()
			board.tiles[idx] = &Tile{interest_queue_pos: -1, xy: xy}
			if x == 1 {
				board.tiles[idx].west_restriction = Border
				board.IncRestrictionCount(board.tiles[idx])
			}
			if x == BOARD_SIZE {
				board.tiles[idx].east_restriction = Border
				board.IncRestrictionCount(board.tiles[idx])
			}
			if y == 1 {
				board.tiles[idx].north_restriction = Border
				board.IncRestrictionCount(board.tiles[idx])
			}
			if y == BOARD_SIZE {
				board.tiles[idx].south_restriction = Border
				board.IncRestrictionCount(board.tiles[idx])
			}
		}
	}

	board.currentPiece += 1
	for _, pp := range board.pieceLookup[board.tiles[0].LookupKey()].GetPieces() {
		if pp.piece.sticky {
			board.Place(pp, Coordinate{x: 1, y: 1})
			break
		}
	}

	board.solveStart = time.Now().UnixMilli()
	return board
}

func (b *Board) IsSolved() bool {
	return b.currentPiece == BOARD_SIZE*BOARD_SIZE
}

func BuildLookup(pieces []Piece) [LOOKUP_SIZE]*PiecePlacementLookup {
	pieceLookup := [LOOKUP_SIZE]*PiecePlacementLookup{}
	for _, p := range pieces {
		for _, pp := range p.Rotations() {
			for _, k := range pp.Keys() {
				if pieceLookup[k] == nil {
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

func (b *Board) GetNextCoordinate() Coordinate {
	minChoices := uint8(255)
	bestCoordinate := Coordinate{x: 0, y: 0}

	for i := 0; i < b.interest_queue_len; i++ {
		t := b.interest_queue[i]
		lookupKey := t.LookupKey()
		lookup := b.pieceLookup[lookupKey]
		choices := uint8(0)
		if lookup != nil {
			choices = lookup.count
		}
		if choices < minChoices {
			minChoices = choices
			bestCoordinate = t.xy
		}
	}

	if bestCoordinate.x == 0 {
		panic("No more interesting tiles")
	}

	IterationOrder[b.currentPiece] = bestCoordinate
	return IterationOrder[b.currentPiece]
}

func (b *Board) PlaceNext() bool {
	// get coordinate of next tile
	coord := b.GetNextCoordinate()
	idx := coord.AsIndex()
	tile := b.tiles[idx]
	backtracking_queue := make([]*PiecePlacement, 0, 1)

	// get matching pieces
	matchingPieces := b.pieceLookup[tile.LookupKey()]
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
		fmt.Printf("  Restrictions for this tile: N=%d E=%d S=%d W=%d\n", tile.north_restriction, tile.east_restriction, tile.south_restriction, tile.west_restriction)
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
		fmt.Printf("  Placing piece %d facing %s (%d, %d, %d, %d)\n", nextCandidate.piece.number, nextCandidate.orientation, nextCandidate.north, nextCandidate.east, nextCandidate.south, nextCandidate.west)
	}
	b.Place(nextCandidate, coord)

	if b.currentPiece > b.maxPlacedPieces {
		b.maxPlacedPieces = b.currentPiece
		if LOG_LEVEL <= WARN {
			fmt.Printf("  Achievement unlocked! Placed %d pieces\n", b.maxPlacedPieces)
		}
	}
	return true
}

func (b *Board) Place(pp *PiecePlacement, xy Coordinate) {
	MovementCount += 1
	if MovementCount%10000000 == 0 && LOG_LEVEL >= WARN {
		timeElapsed := time.Now().UnixMilli() - b.solveStart
		movementsPerSecond := float64(MovementCount) / float64(timeElapsed) / 1000
		fmt.Printf("  %dM movements so far, current rate %.2fM movements per second\n", MovementCount/1000000, movementsPerSecond)
	}
	idx := xy.AsIndex()
	tile := b.tiles[idx]
	tile.placed_piece = pp
	if tile.interest_queue_pos != -1 {
		b.RemoveFromInterestQueue(tile)
	}
	pp.MarkUsed()

	// update restrictions on neighboring tiles and check if it doesn't create an unsolvable situation
	if xy.x > 1 {
		left_tile := b.tiles[idx-1]
		if left_tile.east_restriction == 0 {
			left_tile.east_restriction = pp.west
			b.IncRestrictionCount(left_tile)
		}
	}
	if xy.x < BOARD_SIZE {
		right_tile := b.tiles[idx+1]
		if right_tile.west_restriction == 0 {
			right_tile.west_restriction = pp.east
			b.IncRestrictionCount(right_tile)
		}
	}
	if xy.y > 1 {
		top_tile := b.tiles[idx-BOARD_SIZE]
		if top_tile.south_restriction == 0 {
			top_tile.south_restriction = pp.north
			b.IncRestrictionCount(top_tile)
		}
	}
	if xy.y < BOARD_SIZE {
		bottom_tile := b.tiles[idx+BOARD_SIZE]
		if bottom_tile.north_restriction == 0 {
			bottom_tile.north_restriction = pp.south
			b.IncRestrictionCount(bottom_tile)
		}
	}
}

func (b *Board) Unplace(xy Coordinate) {
	MovementCount += 1
	if MovementCount%10000000 == 0 && LOG_LEVEL >= WARN {
		timeElapsed := time.Now().UnixMilli() - b.solveStart
		movementsPerSecond := float64(MovementCount) / float64(timeElapsed) / 1000
		fmt.Printf("  %dM movements so far, current rate %.2fM movements per second\n", MovementCount/1000000, movementsPerSecond)
	}
	idx := xy.AsIndex()
	tile := b.tiles[idx]
	tile.placed_piece.MarkUnused()
	tile.placed_piece = nil
	if tile.restriction_count >= 2 {
		b.AddToInterestQueue(tile)
	}

	// update restrictions on neighboring tiles
	if xy.x > 1 {
		left_tile := b.tiles[idx-1]
		if left_tile.east_restriction != 0 {
			left_tile.east_restriction = 0
			b.DecRestrictionCount(left_tile)
		}
	}
	if xy.x < BOARD_SIZE {
		right_tile := b.tiles[idx+1]
		if right_tile.west_restriction != 0 {
			right_tile.west_restriction = 0
			b.DecRestrictionCount(right_tile)
		}
	}
	if xy.y > 1 {
		top_tile := b.tiles[idx-BOARD_SIZE]
		if top_tile.south_restriction != 0 {
			top_tile.south_restriction = 0
			b.DecRestrictionCount(top_tile)
		}
	}
	if xy.y < BOARD_SIZE {
		bottom_tile := b.tiles[idx+BOARD_SIZE]
		if bottom_tile.north_restriction != 0 {
			bottom_tile.north_restriction = 0
			b.DecRestrictionCount(bottom_tile)
		}
	}
}

func (b *Board) Backtrack() bool {
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
		if b.currentPiece == 1 {
			return false
		} else {
			return b.Backtrack()
		}
	}

	nextCandidate := tile.Dequeue()
	if LOG_LEVEL <= INFO {
		fmt.Printf("  Placing piece %d facing %s\n", nextCandidate.piece.number, nextCandidate.orientation)
	}
	b.Place(nextCandidate, currentCoordinate)
	return true
}

func (xy Coordinate) AsIndex() int {
	return (xy.y-1)*BOARD_SIZE + xy.x - 1
}

func (b *Board) IncRestrictionCount(t *Tile) {
	t.restriction_count += 1
	if t.restriction_count == 2 {
		b.AddToInterestQueue(t)
	}
}

func (b *Board) DecRestrictionCount(t *Tile) {
	t.restriction_count -= 1
	if t.restriction_count == 1 {
		b.RemoveFromInterestQueue(t)
	}
}

func (b *Board) AddToInterestQueue(t *Tile) {
	if t.interest_queue_pos != -1 {
		return
	}
	b.interest_queue[b.interest_queue_len] = t
	t.interest_queue_pos = b.interest_queue_len
	b.interest_queue_len += 1
}

func (b *Board) RemoveFromInterestQueue(t *Tile) {
	b.interest_queue_len -= 1
	if t.interest_queue_pos != b.interest_queue_len {
		other_tile := b.interest_queue[b.interest_queue_len]
		other_tile.interest_queue_pos = t.interest_queue_pos
		b.interest_queue[t.interest_queue_pos] = other_tile
	}
	t.interest_queue_pos = -1
}
