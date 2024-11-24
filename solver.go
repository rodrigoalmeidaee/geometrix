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

func main() {
	var (
		mode          string
		sortAlgorithm string
	)

	flag.StringVar(&mode, "mode", "solve", "Mode to run the program in")
	flag.StringVar(&sortAlgorithm, "sort", "none", "Sort algorithm to use")
	flag.Parse()

	if mode == "solve" {
		board := Solve(sortAlgorithm)
		if board != nil {
			fmt.Fprintf(os.Stderr, "Solved in %d movements!\n", piece.MovementCount)
			fmt.Printf("%s", *board)
		} else {
			fmt.Fprintf(os.Stderr, "No solution found after %d movements.\n", piece.MovementCount)
		}
	} else if mode == "profile" {
		Profile(10000, sortAlgorithm)
	}
}

func Profile(numAttempts int, sortAlgorithm string) {
	movementCounts := make([]int, numAttempts)
	timings := make([]float64, numAttempts)
	solutions := make(map[string]int)

	for i := 0; i < numAttempts; i++ {
		start := time.Now()
		piece.MovementCount = 0
		board := Solve(sortAlgorithm)
		if board == nil {
			fmt.Fprintf(os.Stderr, "No solution found after %d movements.\n", piece.MovementCount)
			return
		}
		movementCounts[i] = piece.MovementCount
		timings[i] = float64(time.Since(start).Microseconds()) / 1000.0
		if solutions[board.String()] == 0 {
			solutions[board.String()] = len(solutions) + 1
		}
	}

	fmt.Printf("Movements: min=%d, max=%d, avg=%d\n", slices.Min(movementCounts), slices.Max(movementCounts), Avg(movementCounts))
	fmt.Printf("Timings: min=%.2f, max=%.2f, avg=%.2f\n", slices.Min(timings), slices.Max(timings), FloatAvg(timings))
	fmt.Printf("Distinct solutions: %d\n", len(solutions))

	for html, index := range solutions {
		f, _ := os.Create(fmt.Sprintf("pyutils/output-%d.html", index))
		f.WriteString(html)
		f.Close()
	}
}

func Solve(sortAlgorithm string) *piece.Board {
	pieces := piece.GetPieces()
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
