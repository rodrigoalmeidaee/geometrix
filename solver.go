package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"rsalmeidafl/geometrix/piece"
	"slices"
	"time"
)

func GetPieces() []piece.Piece {
	piece.ResetPieceNumber()

	var pieces = []piece.Piece{
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
	return pieces
}

func main() {
	var (
		mode string
	)

	flag.StringVar(&mode, "mode", "solve", "Mode to run the program in")
	flag.Parse()

	if mode == "solve" {
		board := Solve()
		if board != nil {
			fmt.Fprintf(os.Stderr, "Solved in %d movements!\n", piece.MovementCount)
			fmt.Printf("%s", *board)
		} else {
			fmt.Fprintf(os.Stderr, "No solution found after %d movements.\n", piece.MovementCount)
		}
	} else if mode == "profile" {
		Profile(1000)
	}
}

func Profile(numAttempts int) {
	movementCounts := make([]int, numAttempts)
	timings := make([]float64, numAttempts)

	for i := 0; i < numAttempts; i++ {
		start := time.Now()
		piece.MovementCount = 0
		if Solve() == nil {
			fmt.Fprintf(os.Stderr, "No solution found after %d movements.\n", piece.MovementCount)
			return
		}
		movementCounts[i] = piece.MovementCount
		timings[i] = float64(time.Since(start).Microseconds()) / 1000.0
	}

	fmt.Printf("Movements: min=%d, max=%d, avg=%d\n", slices.Min(movementCounts), slices.Max(movementCounts), Avg(movementCounts))
	fmt.Printf("Timings: min=%.2f, max=%.2f, avg=%.2f\n", slices.Min(timings), slices.Max(timings), FloatAvg(timings))
}

func Solve() *piece.Board {
	pieces := GetPieces()
	perm := rand.Perm(len(pieces))
	shuffledPieces := make([]piece.Piece, len(pieces))

	for i, v := range perm {
		shuffledPieces[v] = pieces[i]
	}

	board := piece.NewBoard(shuffledPieces)

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

func Avg(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum / len(nums)
}

func FloatAvg(nums []float64) float64 {
	sum := float64(0)
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums))
}
