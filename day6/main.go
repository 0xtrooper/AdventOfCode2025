package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:embed input.txt
var rawInput string

var rawLines []string

type Operation struct {
	Offset int
	Length int
	Fn     func(int, int) int
}

var operations []Operation

func parseInput() {
	scanner := bufio.NewScanner(strings.NewReader(rawInput))

	// read all lines
	for scanner.Scan() {
		rawLines = append(rawLines, scanner.Text())
	}

	// parse math operations
	mathOperationsLine := rawLines[len(rawLines)-1]

	index := 0
ParseOperations:
	for index < len(mathOperationsLine) {
		nextOperation := Operation{Offset: index, Length: 1}

		// parse operation
		switch mathOperationsLine[index] {
		case '+':
			nextOperation.Fn = func(a, b int) int { return a + b }
		case '*':
			nextOperation.Fn = func(a, b int) int { return a * b }
		default:
			panic(fmt.Sprintf("unknown operation: %c", mathOperationsLine[index]))
		}

		// find next operator or end of line
		for {
			if index+nextOperation.Length+1 >= len(mathOperationsLine) {
				nextOperation.Length += 1
				operations = append(operations, nextOperation)
				break ParseOperations
			}
			next := mathOperationsLine[index+nextOperation.Length+1]
			if next == '+' || next == '*' {
				index += nextOperation.Length + 1
				operations = append(operations, nextOperation)
				continue ParseOperations
			}
			nextOperation.Length++
		}
	}
}

func partOne() int {
	getNumber := func(index, offset, length int) int {
		numStr := rawLines[index][offset : offset+length]
		numStr = strings.TrimSpace(numStr)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic(fmt.Sprintf("%d - [%d:%d] invalid number: '%s'", index, offset, length, numStr))
		}
		return num
	}

	total := 0
	for _, operation := range operations {
		operationRes := getNumber(0, operation.Offset, operation.Length)
		for i := 1; i < len(rawLines)-1; i++ {
			operationRes = operation.Fn(operationRes, getNumber(i, operation.Offset, operation.Length))
		}
		total += operationRes
	}
	return total
}

func partTwo() int {
	// read number top to bottom at given index and offset
	getNumber := func(offset, index int) int {
		numStr := string(rawLines[0][offset+index])
		for i := 1; i < len(rawLines)-1; i++ {
			numStr += string(rawLines[i][offset+index])
		}
		numStr = strings.TrimSpace(numStr)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic(fmt.Sprintf("index %d offset %d invalid number: '%s'", index, offset, numStr))
		}
		return num
	}

	total := 0
	for _, operation := range operations {
		operationRes := getNumber(operation.Offset, 0)
		for i := 1; i < operation.Length; i++ {
			operationRes = operation.Fn(operationRes, getNumber(operation.Offset, i))
		}
		total += operationRes
	}
	return total
}

func run(fn func() int) int {
	startTime := time.Now()
	result := fn()
	elapsed := time.Since(startTime)
	fmt.Printf("Evaluation time: %s  - ", elapsed)
	return result
}

func main() {
	parseInput()
	fmt.Printf("Part One: %d\n", run(partOne))
	fmt.Printf("Part Two: %d\n", run(partTwo))
}
