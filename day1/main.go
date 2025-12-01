package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var rawInput string

const (
	maxCount      = 100
	startPosition = 50
)

// parseInput processes the embedded input and returns a slice of integers representing the turn directions.
// direction > 0 means turn right, < 0 means turn left.
func parseInput() ([]int, error) {
	// Clean up any carriage return characters
	instructions := strings.Split(strings.ReplaceAll(rawInput, "\r", ""), "\n")
	directions := make([]int, len(instructions))

	for i, instr := range instructions {

		if strings.HasPrefix(instr, "R") {
			directions[i] = 1
		} else if strings.HasPrefix(instr, "L") {
			directions[i] = -1
		} else {
			return nil, fmt.Errorf("invalid instruction: %s", instr)
		}

		value, err := strconv.Atoi(instr[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid number in instruction: '%s' ['%s'] - %s", instr, instr[1:], err.Error())
		}

		directions[i] *= value
	}
	return directions, nil
}

// executeTurnsCountZeroHits processes the list of turn directions and returns how often the dial hit 0
func executeTurnsCountZeroHits(turnDirections []int) int {
	position := startPosition
	hitZeroCount := 0

	for _, turn := range turnDirections {
		position = (position + turn) % maxCount
		if position < 0 {
			position += maxCount
		}

		if position == 0 {
			hitZeroCount++
		}
	}
	return hitZeroCount
}

func executeTurnsCrossZero(turnDirections []int) int {
	position := startPosition
	crossedZeroCount := 0

	for _, turn := range turnDirections {
		nextPosition := position + turn

		if nextPosition <= 0 && position != 0 {
			crossedZeroCount++
		}

		if nextPosition > 0 {
			crossedZeroCount += nextPosition / maxCount
		} else {
			crossedZeroCount += (-nextPosition) / maxCount
		}

		position = ((nextPosition % maxCount) + maxCount) % maxCount
	}

	return crossedZeroCount
}

func main() {
	turnDirections, err := parseInput()
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	hitZeroCount := executeTurnsCountZeroHits(turnDirections)
	fmt.Printf("[Part 1] The dial hit 0 a total of %d times.\n", hitZeroCount)

	crossedZeroCount := executeTurnsCrossZero(turnDirections)
	fmt.Printf("[Part 2] The dial crossed 0 a total of %d times.\n", crossedZeroCount)
}
