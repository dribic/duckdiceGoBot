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
	"math"
	"os"
	"time"
)

func LabouchereSpec(startBet, startBalance, targetBalance float64, mode, high bool, apiKey, hash, curr string) float64 {
	// Initialise the sequence table
	seqTable := make([]float64, 10)
	for i := range seqTable {
		seqTable[i] = startBet
	}

	// Initialise other variables
	currentBalance := startBalance
	var totalBetAmount float64
	var multiplier float64 = 1.2045
	var victories, loses uint16
	safety := false

	fmt.Println("Win multiplier is", multiplier)

	// Main loop
	for currentBalance < targetBalance {
		var currentBet, bet1, bet2 float64
		if len(seqTable) == 1 {
			bet1 = seqTable[0]
			currentBet = bet1
			currentBet = math.Round(currentBet*100000) / 100000
		} else {
			bet1, bet2 = seqTable[0], seqTable[len(seqTable)-1]
			currentBet = bet1 + bet2
			currentBet = math.Round(currentBet*100000) / 100000
		}

		if currentBalance < currentBet {
			fmt.Println("Balance too low!")
			os.Exit(0)
		}

		// Rebuilding the betting sequence, if the bet gets too big
		if seqTable[0] >= 3*startBet || seqTable[len(seqTable)-1] >= 10*startBet {
			fmt.Println("Bets have became too large.")
			fmt.Println("Constructing new betting sequence...")
			seqTable = LabSafety(startBet, seqTable)
			fmt.Println("New betting sequence:", seqTable)
			time.Sleep(time.Second * 3)
		}

		// Lowering the currentBet if larger than difference to targetBalance
		for currentBet > targetBalance-currentBalance && currentBet > startBet {
			safety = true
			fmt.Println("Lowering", currentBet, "by", startBet, "as a precaution.")
			currentBet -= startBet
			currentBet = math.Round(currentBet*100000) / 100000
		}

		totalBetAmount += currentBet

		// Print what the next bet is
		fmt.Printf("Next bet: %.6f (%.6f + %.6f)\n", currentBet, bet1, bet2)

		// Turn bet into a string, because DDs API uses string
		amount := fmt.Sprintf("%.6f", currentBet)

		// Making a bet
		result, _ := PlaceABetSpec(apiKey, amount, "44", curr, hash, mode, high, true)

		if result {
			if safety {
				fmt.Println("Success!✅")
				victories++
				currentBalance += currentBet * multiplier
				if seqTable[0] == currentBet {
					seqTable = seqTable[1:]
				} else {
					var ticker float64 = 0.0
					for ticker < currentBet {
						ticker += startBet
						seqTable[len(seqTable)-1] -= startBet
						if seqTable[len(seqTable)-1] < startBet {
							seqTable = seqTable[:len(seqTable)-1]
						}
					}
				}
			} else {
				fmt.Println("Success!✅")
				victories++
				currentBalance += currentBet * multiplier
				seqTable = seqTable[1 : len(seqTable)-1]
			}
		} else {
			fmt.Println("Failure!☯ ")
			loses++
			currentBalance -= currentBet
			currentBet = math.Round(currentBet*100000) / 100000
			seqTable = append(seqTable, currentBet)
		}

		// Print progress
		fmt.Printf("Current balance: %.6f, Victories: %d, Loses: %d, Total gambled amount: %6f\n", currentBalance,
			victories, loses, totalBetAmount)
		fmt.Println("Betting sequence table:", seqTable)
	}

	return currentBalance
}
