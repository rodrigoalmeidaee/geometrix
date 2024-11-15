package main

import (
	"fmt"
	"math/rand"
	"os"
	"rsalmeidafl/geometrix/piece"

	"github.com/emirpasic/gods/lists/arraylist"
)

var pieces = [36]piece.Piece{
	piece.New(piece.BlueStar, piece.YellowGreenTriangle, piece.Border, piece.RedYellowTriangle),
	piece.New(piece.BlueStar, piece.PinkTrident, piece.BlueStar, piece.YellowGreenTriangle),
	piece.New(piece.PinkTrident, piece.BlueStar, piece.YellowCircle, piece.Border),
	piece.New(piece.PinkCircle, piece.RedYellowTriangle, piece.Border, piece.RedYellowTriangle),
	piece.New(piece.RedYellowTriangle, piece.BlueStar, piece.PinkCircle, piece.YellowGreenTriangle),
	piece.New(piece.PinkCircle, piece.PinkTrident, piece.RedTrident, piece.Border),
	piece.New(piece.Border, piece.RedTrident, piece.YellowCircle, piece.YellowGreenTriangle),
	piece.New(piece.RedYellowTriangle, piece.PinkCircle, piece.YellowCircle, piece.BlueStar),
	piece.New(piece.YellowCircle, piece.RedTrident, piece.PinkCircle, piece.Border),
	piece.New(piece.YellowCircle, piece.Border, piece.Border, piece.RedYellowTriangle),
	piece.New(piece.YellowGreenTriangle, piece.PinkCircle, piece.RedYellowTriangle, piece.PinkTrident),
	piece.New(piece.PinkCircle, piece.Border, piece.YellowCircle, piece.BlueStar),
	piece.New(piece.YellowGreenTriangle, piece.RedTrident, piece.Border, piece.Border),
	piece.New(piece.RedTrident, piece.PinkTrident, piece.BlueStar, piece.YellowGreenTriangle),
	piece.New(piece.YellowCircle, piece.Border, piece.RedTrident, piece.YellowGreenTriangle),
	piece.New(piece.Border, piece.YellowGreenTriangle, piece.YellowCircle, piece.RedYellowTriangle),
	piece.New(piece.YellowCircle, piece.PinkCircle, piece.YellowGreenTriangle, piece.Border),
	piece.New(piece.RedTrident, piece.YellowCircle, piece.RedYellowTriangle, piece.YellowGreenTriangle),
	piece.New(piece.BlueStar, piece.RedYellowTriangle, piece.Border, piece.PinkTrident),
	piece.New(piece.YellowGreenTriangle, piece.YellowCircle, piece.RedTrident, piece.RedYellowTriangle),
	piece.New(piece.PinkTrident, piece.Border, piece.YellowCircle, piece.PinkCircle),
	piece.New(piece.BlueStar, piece.YellowGreenTriangle, piece.RedTrident, piece.PinkCircle),
	piece.New(piece.PinkCircle, piece.RedTrident, piece.PinkTrident, piece.Border),
	piece.New(piece.PinkCircle, piece.RedTrident, piece.YellowCircle, piece.YellowGreenTriangle),
	piece.New(piece.RedYellowTriangle, piece.RedTrident, piece.Border, piece.Border),
	piece.New(piece.PinkTrident, piece.RedTrident, piece.YellowGreenTriangle, piece.PinkCircle),
	piece.New(piece.RedTrident, piece.PinkCircle, piece.YellowCircle, piece.Border),
	piece.New(piece.PinkCircle, piece.PinkTrident, piece.RedTrident, piece.PinkTrident),
	piece.New(piece.BlueStar, piece.PinkCircle, piece.PinkTrident, piece.PinkCircle),
	piece.New(piece.YellowCircle, piece.PinkCircle, piece.Border, piece.PinkTrident),
	piece.New(piece.PinkCircle, piece.YellowCircle, piece.BlueStar, piece.YellowCircle),
	piece.New(piece.PinkTrident, piece.RedYellowTriangle, piece.Border, piece.PinkCircle),
	piece.New(piece.YellowGreenTriangle, piece.PinkTrident, piece.Border, piece.Border),
	piece.New(piece.RedTrident, piece.YellowCircle, piece.BlueStar, piece.YellowCircle),
	piece.New(piece.BlueStar, piece.PinkTrident, piece.BlueStar, piece.YellowGreenTriangle),
	piece.New(piece.RedYellowTriangle, piece.YellowCircle, piece.BlueStar, piece.YellowGreenTriangle),
}

var pieceLookup map[int]*arraylist.List

func main() {
	perm := rand.Perm(36)
	shuffledPieces := [36]piece.Piece{}

	for i, v := range perm {
		shuffledPieces[v] = pieces[i]
	}

	pieceLookup = make(map[int]*arraylist.List)
	for _, p := range shuffledPieces {
		for _, pp := range p.Rotations() {
			for _, k := range pp.Keys() {
				if _, ok := pieceLookup[k]; !ok {
					pieceLookup[k] = arraylist.New()
				}
				pieceLookup[k].Add(&pp)
			}
		}
	}

	board := piece.NewBoard()

	for {
		if board.PlaceNext(pieceLookup) {
			if board.IsSolved() {
				fmt.Fprintf(os.Stderr, "Solved in %d movements!\n", piece.MovementCount)
				fmt.Printf("%s", board)
				break
			}
		} else {
			if board.Backtrack() {
				continue
			} else {
				fmt.Fprintf(os.Stderr, "No solution found after %d movements.\n", piece.MovementCount)
				break
			}
		}
	}
}
