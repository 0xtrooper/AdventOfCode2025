package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	CONSOLE_RESET = "\u001B[0m"
	CONSOLE_RED   = "\u001B[31m"
)

//go:embed input.txt
var rawInput string

var batterysPerBank = 0

type BatteryBank struct {
	CellVoltate         []string
	CellIndexesToTurnOn []int
}

func NewBatteryBank(s string) (b BatteryBank, err error) {
	return BatteryBank{
		CellVoltate: strings.Split(s, ""),
	}, nil
}

func (b BatteryBank) String() string {
	builder := strings.Builder{}
	for i, d := range b.CellVoltate {
		if slices.Contains(b.CellIndexesToTurnOn, i) {
			builder.WriteString(CONSOLE_RED)
			builder.WriteString(d)
			builder.WriteString(CONSOLE_RESET)
		} else {
			builder.WriteString(d)
		}
	}
	return builder.String()
}

func (b *BatteryBank) GetMaxJoltage(cellCount int) (int, error) {
	if len(b.CellIndexesToTurnOn) != cellCount {
		if err := b.computeMaxJoltage(cellCount); err != nil {
			return 0, err
		}
	}

	sum := 0
	for _, i := range b.CellIndexesToTurnOn {
		sum = sum*10 + int(b.CellVoltate[i][0]-'0')
	}
	return sum, nil
}

func (b *BatteryBank) computeMaxJoltage(cellCount int) error {
	if cellCount > len(b.CellVoltate) {
		return fmt.Errorf("requested cell count exceeds length of the battery bank (%d > %d)", cellCount, len(b.CellVoltate))
	}

	// Assume the higest digit is in the first place
	b.CellIndexesToTurnOn = make([]int, cellCount)
	for i := range cellCount {
		b.CellIndexesToTurnOn[i] = i
	}

	// keep track of how many 9's we found, so we do not constantly re-check as they can not get bigger
	ninesFound := 0

	// Iterate over the follwing digits, prioritze the first digit
	// iterate up to N-cellCount, as a last resort, we just return the final 'cellCount' digits
	for i := 1; i < batterysPerBank-1; i++ {
		// Start indexing the starting from the higest remaining index
		offset := max(0, (i+cellCount)-batterysPerBank)

		// if we found a new 9, update the start index
		if b.CellVoltate[ninesFound] == "9" {
			ninesFound++
		}

		// starting at the highest index, check if the current cell volate is higer than the sorted one
		for j := ninesFound; j+offset < cellCount; j++ {
			if b.CellVoltate[i+j] > b.CellVoltate[b.CellIndexesToTurnOn[j+offset]] {
				// the number in higher as the one we have right now, we have to reset the following indexes
				for ; j+offset < cellCount; j++ {
					b.CellIndexesToTurnOn[j+offset] = i + j
				}

				break
			}
		}
	}
	return nil
}

func parseBatteryBanks(input string) (banks []BatteryBank, err error) {
	bankStrs := strings.Split(input, "\n")

	// Assume all banks have the same number of batteries (will validate later)
	batterysPerBank = len(strings.TrimSpace(bankStrs[0]))

	banks = make([]BatteryBank, len(bankStrs))
	for i, bankStr := range bankStrs {
		bankStr = strings.TrimSpace(bankStr)
		if bankStr == "" {
			return nil, fmt.Errorf("empty bank string at line %d", i+1)
		}

		// Verify that the bank string has the correct length
		if len(bankStr) != batterysPerBank {
			return nil, fmt.Errorf("inconsistent battery count at line %d - '%s': expected %d, got %d", i+1, bankStr, batterysPerBank, len(bankStr))
		}

		banks[i], err = NewBatteryBank(bankStr)
		if err != nil {
			return nil, fmt.Errorf("invalid bank string at line %d: %v", i+1, err)
		}
	}

	return banks, nil
}

func main() {
	startTime := time.Now()
	batteryBanks, err := parseBatteryBanks(rawInput)
	if err != nil {
		fmt.Println("Error parsing battery banks:", err)
		return
	}

	sum2 := 0
	sum12 := 0
	for i, bank := range batteryBanks {
		maxJoltage2, err := bank.GetMaxJoltage(2)
		if err != nil {
			panic(err)
		}

		maxJoltage12, err := bank.GetMaxJoltage(12)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[Bank %d]: %s - Max Joltage: %d\n", i+1, bank.String(), maxJoltage12)

		sum2 += maxJoltage2
		sum12 += maxJoltage12
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nSum of max joltages for 2 cells: %d; for 12 cells: %d - elapsed time: %s\n", sum2, sum12, elapsed)
}
