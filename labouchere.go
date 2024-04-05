/*
 Copyright (C) 2024 Dejan Ribič

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
)

func Labouchere(startBet, startBalance, targetBalance float64, fMode, high bool, apiKey, curr string) float64 {
	// Initialise the sequence table
	seqTable := make([]float64, 10)
	for i := range seqTable {
		seqTable[i] = startBet
	}

	// Initialise other variables
	currentBalance := startBalance
	var totalBetAmount float64
	var multiplier float64 = 1.25
	var victories, loses uint16

	// Lowering multiplier if fMode == true
	if fMode {
		multiplier = 1.2045
	}

	fmt.Println("Win multiplier is", multiplier)

	// Main loop
	for currentBalance < targetBalance {
		var currentBet, bet1, bet2 float64
		if len(seqTable) == 1 {
			bet1 = seqTable[0]
			currentBet = bet1
		} else {
			bet1, bet2 = seqTable[0], seqTable[len(seqTable)-1]
			currentBet = bet1 + bet2
		}

		if currentBalance < currentBet {
			fmt.Println("Balance too low!")
			os.Exit(0)
		}

		totalBetAmount += currentBet

		// Print what the next bet is
		fmt.Printf("Next bet: %.6f (%.6f + %.6f)\n", currentBet, bet1, bet2)

		// Turn bet into a string, because DDs API uses string
		amount := fmt.Sprintf("%.6f", currentBet)

		// Making a bet
		result := PlaceABet(apiKey, amount, "44", curr, fMode, high)

		if result {
			fmt.Println("Success!✅")
			victories++
			currentBalance += currentBet * multiplier
			seqTable = seqTable[1 : len(seqTable)-1]
		} else {
			fmt.Println("Failure!☯ ")
			loses++
			currentBalance -= currentBet
			seqTable = append(seqTable, bet1+bet2)
		}

		// Print progress
		fmt.Printf("Current balance: %.6f, Victories: %d, Loses: %d, Total gambled amount: %6f\n", currentBalance,
			victories, loses, totalBetAmount)
		fmt.Println("Betting sequence table:", seqTable)
	}

	return currentBalance
}
