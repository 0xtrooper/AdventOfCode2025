package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed input.txt
var rawInput string

const (
	START_SYMBOL    = "S"
	SPLITTER_SYMBOL = "^"
	EMPTY_SYMBOL    = "."
	BEAM_SYMBOL     = "|"

	EMPTY_RUNE    = '.'
	SPLITTER_RUNE = '^'
)

func partOne() (total int) {
	// processBeam processes the current line based on the previous line and returns the updated current line.
	processBeam := func(previousLineStr, currentLineStr string) (string, int) {
		previousLineChars := strings.Split(previousLineStr, "")
		currentLineChars := strings.Split(currentLineStr, "")

		if len(previousLineChars) != len(currentLineChars) {
			msg := fmt.Sprintf("Lines must be of equal length (%d != %d)", len(previousLineChars), len(currentLineChars))
			panic(msg)
		}

		splitCount := 0
		for i := range currentLineChars {
			prevChar := previousLineChars[i]
			currChar := currentLineChars[i]

			if prevChar == BEAM_SYMBOL || prevChar == START_SYMBOL {
				switch currChar {
				case EMPTY_SYMBOL:
					currentLineChars[i] = BEAM_SYMBOL
				case SPLITTER_SYMBOL:
					splitCount++
					// Split beam to left and right
					if i > 0 && currentLineChars[i-1] == EMPTY_SYMBOL {
						currentLineChars[i-1] = BEAM_SYMBOL
					}
					if i < len(currentLineChars)-1 && currentLineChars[i+1] == EMPTY_SYMBOL {
						currentLineChars[i+1] = BEAM_SYMBOL
					}
				}
			}
		}

		return strings.Join(currentLineChars, ""), splitCount
	}

	scanner := bufio.NewScanner(strings.NewReader(rawInput))
	scanner.Scan()
	line := scanner.Text()
	for scanner.Scan() {
		var splitCount int
		line, splitCount = processBeam(line, scanner.Text())
		total += splitCount
	}

	return total
}

type keyType struct {
	lineIndex uint8
	beamIndex uint8
}

var cache map[keyType]int

func eval(s [][]bool, lineIndex, beamIndex uint8) int {
	if lineIndex >= uint8(len(s)) {
		return 1
	}

	key := keyType{
		lineIndex: lineIndex,
		beamIndex: beamIndex,
	}

	if res, ok := cache[key]; ok {
		return res
	}

	var res int
	if s[lineIndex][beamIndex] {
		res = eval(s, lineIndex+1, beamIndex-1) + eval(s, lineIndex+1, beamIndex+1)
	} else {
		res = eval(s, lineIndex+1, beamIndex)
	}

	cache[key] = res

	return res
}

func partTwo() int {
	scanner := bufio.NewScanner(strings.NewReader(rawInput))
	scanner.Scan()
	line := scanner.Text()

	// look for start index
	beamIndex := uint8(strings.Index(line, START_SYMBOL))

	splitters := [][]bool{}
	for scanner.Scan() {
		line := scanner.Text()

		splitterLine := make([]bool, len(line))
		for i, elem := range line {
			splitterLine[i] = (elem == SPLITTER_RUNE)
		}
		splitters = append(splitters, splitterLine)
	}

	cache = make(map[keyType]int)
	return eval(splitters, 0, beamIndex)
}

func run(fn func() int) int {
	startTime := time.Now()
	result := fn()
	elapsed := time.Since(startTime)
	fmt.Printf("Evaluation time: %s  - ", elapsed)
	return result
}

func main() {
	fmt.Printf("Part One: %d\n", run(partOne))
	fmt.Printf("Part Two: %d\n", run(partTwo))
}
