package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"rsalmeidafl/geometrix/piece"
)

func main() {
	var (
		sortAlgorithm     string
		nextTilePlacement string
		targetGeneration  int
	)

	flag.StringVar(&sortAlgorithm, "sort", "none", "Sort algorithm to use")
	flag.IntVar(&targetGeneration, "gen", piece.BOARD_SIZE, "Solve only up to this generation")
	flag.StringVar(&nextTilePlacement, "nextpiece", "least-options", "How to pick the next tile to place")
	flag.Parse()

	board := Solve(sortAlgorithm, targetGeneration, nextTilePlacement)
	if board != nil {
		fmt.Fprintf(os.Stderr, "Solved in %d movements!\n", piece.MovementCount)
		fmt.Printf("%s", *board)
	} else {
		fmt.Fprintf(os.Stderr, "No solution found after %d movements.\n", piece.MovementCount)
	}
}

func Solve(sortAlgorithm string, targetGeneration int, nextTilePlacement string) *piece.Board {
	pieces := piece.GetPieces()
	perm := rand.Perm(len(pieces))
	shuffledPieces := make([]piece.Piece, len(pieces))

	for i, v := range perm {
		shuffledPieces[v] = pieces[i]
	}

	board := piece.NewBoard(shuffledPieces, targetGeneration, nextTilePlacement)

	for {
		if board.PlaceNext() {
			if board.IsSolved() {
				return &board
			}
		} else {
			if !board.Backtrack() {
				return nil
			}
		}
	}
}

func Avg(nums []int64) int64 {
	sum := int64(0)
	for _, n := range nums {
		sum += n
	}
	return sum / int64(len(nums))
}

func FloatAvg(nums []float64) float64 {
	sum := float64(0)
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums))
}
