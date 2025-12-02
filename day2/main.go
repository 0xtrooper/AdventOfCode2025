package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var rawInput string

type Range struct {
	Start int
	End   int
}

func parseRanges(input string) (ranges []Range, err error) {
	rangesStrs := strings.Split(input, ",")
	ranges = make([]Range, len(rangesStrs))
	for i, rStr := range rangesStrs {
		parts := strings.Split(rStr, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range: %s", rStr)
		}
		ranges[i].Start, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start in range '%s': %s", rStr, err.Error())
		}
		ranges[i].End, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid end in range '%s': %s", rStr, err.Error())
		}
	}
	return ranges, nil
}

func filterRangePartOne(r Range) (wrongIds []int, sum int) {
	for id := r.Start; id <= r.End; id++ {
		idStr := strconv.Itoa(id)
		n := len(idStr)

		// uneven length words can not inlude patterns
		if n%2 == 0 && idStr[0:n/2] == idStr[n/2:n] {
			wrongIds = append(wrongIds, id)
			sum += id
		}
	}
	return wrongIds, sum
}

func filterRangePartTwo(r Range) (wrongIds []int, sum int) {
rangeSearch:
	for id := r.Start; id <= r.End; id++ {
		idStr := strconv.Itoa(id)
		n := len(idStr)

	patternSearch:
		for patternLength := 1; patternLength <= n/2; patternLength++ {
			// patternLength must be a divisor of n
			if n%patternLength != 0 {
				continue
			}

			// compare all following segments to the first one
			pattern := idStr[0:patternLength]
			for i := patternLength; i < n; i += patternLength {
				// segment does not match pattern, try next pattern length
				if pattern != idStr[i:i+patternLength] {
					continue patternSearch
				}
			}

			// found repeating pattern
			wrongIds = append(wrongIds, id)
			sum += id

			// skip to next id
			continue rangeSearch
		}
	}
	return wrongIds, sum

}

func main() {
	fmt.Printf("input: %s", rawInput)
	ranges, err := parseRanges(strings.TrimSpace(rawInput))
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed ranges: ")
	sumPartOne := 0
	sumPartTwo := 0
	for _, r := range ranges {
		_, sumRangePartOne := filterRangePartOne(r)
		// fmt.Printf("  Range: %+v - Wrong IDs: %v\n", r, wrongIds)
		sumPartOne += sumRangePartOne

		_, sumRangePartTwo := filterRangePartTwo(r)
		// fmt.Printf("  Range: %+v - Wrong IDs: %v\n", r, wrongIds)
		sumPartTwo += sumRangePartTwo
	}
	fmt.Printf("Total sum of wrong IDs (Part One): %d\n", sumPartOne)
	fmt.Printf("Total sum of wrong IDs (Part Two): %d\n", sumPartTwo)
}
