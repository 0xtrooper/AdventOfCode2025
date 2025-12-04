package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	_ "embed"
)

//go:embed input.txt
var rawInput string

const PAPER_ROLE = '@'

var diagram [][]uint8

func parseDiagram() {
	scanner := bufio.NewScanner(strings.NewReader(rawInput))
	diagram = [][]uint8{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		row := make([]uint8, len(line))
		for i, ch := range line {
			if ch == PAPER_ROLE {
				row[i] = 1
			}
		}
		diagram = append(diagram, row)
	}
}

// evaluateCell evaluates the cell at (row, col); it returns true if the sum of adjacent cells is < 4
func evaluateCell(row, col int) bool {
	sum := 9
	for r := max(0, row-1); r <= min(len(diagram)-1, row+1); r++ {
		for c := max(0, col-1); c <= min(len(diagram[0])-1, col+1); c++ {
			sum -= int(diagram[r][c])
		}
	}
	return diagram[row][col] == 1 && sum > 4 // additional check to ensure we are evaluating a paper cell (normally handled in the caller)
}

func countMovableCells() int {
	count := 0
	for r, row := range diagram {
		for c := range row {
			if diagram[r][c] == 1 && evaluateCell(r, c) {
				count++
			}
		}
	}
	return count
}

func countAndMoveCells() int {
	count := 0

	tempDiagram := make([][]uint8, len(diagram))
	for r := range diagram {
		tempDiagram[r] = make([]uint8, len(diagram[r]))
		copy(tempDiagram[r], diagram[r])
	}

	for {
		// Copy state
		oldCount := count
		for r, row := range diagram {
			for c := range row {
				if diagram[r][c] == 1 && evaluateCell(r, c) {
					diagram[r][c] = 0
					count++
				}
			}
		}

		if count == oldCount {
			break
		}
	}

	// Restore diagram state
	diagram = tempDiagram

	return count
}

func run(fn func() int) int {
	startTime := time.Now()
	result := fn()
	elapsed := time.Since(startTime)
	fmt.Printf("Evaluation time: %s  - ", elapsed)
	return result
}

func main() {
	parseDiagram()
	fmt.Printf("Number of movable paper cells: %d\n", run(countMovableCells))
	fmt.Printf("Number of moved paper cells: %d\n", run(countAndMoveCells))
}
