package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed input.txt
var rawInput string

type databaseEntry struct {
	id      int
	isFresh bool
}

var database []databaseEntry

func addRangeToDatabase(startIdStr, endIdStr string) {
	startId, err := strconv.Atoi(startIdStr)
	if err != nil {
		panic(fmt.Sprintf("invalid startId: %s", startIdStr))
	}

	endId, err := strconv.Atoi(endIdStr)
	if err != nil {
		panic(fmt.Sprintf("invalid endId: %s", endIdStr))
	}
	endId += 1 // make endId exclusive

	// Find first entry >= startId
	idxStart := sort.Search(len(database), func(i int) bool {
		return database[i].id >= startId
	})

	// Find first entry >= endId
	idxEnd := sort.Search(len(database), func(i int) bool {
		return database[i].id >= endId
	})

	// If idxStart > 0, we look at the previous entry.
	isFreshBefore := idxStart > 0 && database[idxStart-1].isFresh

	isFreshAfter := false
	if idxEnd < len(database) && database[idxEnd].id == endId {
		// exact match
		isFreshAfter = database[idxEnd].isFresh
	} else if idxEnd > 0 {
		// look at previous entry
		isFreshAfter = database[idxEnd-1].isFresh
	}

	var newMarkers []databaseEntry

	// add fresh market if the state before is NOT Fresh
	if !isFreshBefore {
		newMarkers = append(newMarkers, databaseEntry{id: startId, isFresh: true})
	}

	// add un-fresh marker if the state after is NOT Fresh
	if !isFreshAfter {
		newMarkers = append(newMarkers, databaseEntry{id: endId, isFresh: false})
	}

	// We remove everything from idxStart up to idxEnd. Add one if we have an exact match at endId.
	replaceUntil := idxEnd
	if idxEnd < len(database) && database[idxEnd].id == endId {
		replaceUntil++
	}

	database = slices.Replace(database, idxStart, replaceUntil, newMarkers...)
}

func setupDatabase(scanner *bufio.Scanner) {
	database = []databaseEntry{}

	// add database entries until we find an empty line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}

		// entry line format: "int-int"
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic(fmt.Sprintf("wrong length on line: %s", line))
		}

		addRangeToDatabase(parts[0], parts[1])
	}
}

func countFreshIngriedients(scanner *bufio.Scanner) int {
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		itemId, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprintf("invalid item id: %s", line))
		}

		// binary search in database
		index := sort.Search(len(database), func(i int) bool {
			return database[i].id > itemId
		})

		if index > 0 && database[index-1].isFresh {
			count++
		}
	}

	return count
}

func countTotalFreshIds() int {
	count := 0
	for i := 0; i < len(database)-1; i++ {
		if database[i].isFresh {
			count += database[i+1].id - database[i].id
		}
	}
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
	scanner := bufio.NewScanner(strings.NewReader(rawInput))

	setupDatabase(scanner)

	fmt.Printf("Part One: %d\n", run(func() int {
		return countFreshIngriedients(scanner)
	}))

	fmt.Printf("Part Two: %d\n", run(countTotalFreshIds))

}
