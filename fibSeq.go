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
	"strings"
)

func addNextFib(first, second int, inTable []int) []int {
	newTable := inTable

	newTable = append(newTable, first+second)

	return newTable
}

func FibBetting(baseBet, startBalance, endBalance, limitLoss float64, fMode, high bool, apiKey, curr string) float64 {
	const first int = 0
	const second int = 1
	var betTable []int = []int{first, second}
	balance := startBalance
	var multiplier float64 = 1.25
	var betNumber uint = 1

	// Lowering multiplier if fMode == true
	if fMode {
		multiplier = 1.2045
	}

	fmt.Println("Win multiplier is", multiplier)

	fmt.Println("Current balance is", balance)
	fmt.Println("Your base bet is", baseBet)
	if limitLoss == 0 {
		fmt.Println("Warning you can lose your entire balance!")
	}

	var choice string
	fmt.Print("Are you sure you wanna run fibonacci sequence <yes/no>? ")
	fmt.Scan(&choice)
	choice = strings.ToUpper(choice)
	if choice == "NO" || choice == "N" {
		return balance
	}

	for balance < endBalance && balance > limitLoss {
		fmt.Println(betTable)
		for len(betTable) > 1 {
			fmt.Printf("%d. bet:\n", betNumber)
			betNumber++
			currentBet := float64(betTable[len(betTable)-1]) * baseBet
			currentBet = math.Round(currentBet*100000) / 100000
			fmt.Println("Next bet is", currentBet)
			if currentBet > balance {
				fmt.Println("Balance too small!!!")
				return balance
			} else if balance-currentBet < limitLoss {
				fmt.Println("Your loss limit has been reached.")
				fmt.Println("Stopping.")
				return balance
			}

			// Placing a bet
			result := PlaceABet(apiKey, fmt.Sprintf("%.8f", currentBet), "44", curr, fMode, high)
			if result {
				fmt.Println("Success!✅")
				balance += currentBet * multiplier
				fmt.Println("Current balance is", balance)
				if balance >= endBalance {
					break
				}
				if len(betTable) == 2 {
					break
				} else {
					betTable = betTable[:len(betTable)-2]
					fmt.Println(betTable)
				}
			} else {
				fmt.Println("Failure!☯ ")
				balance -= currentBet
				if balance < limitLoss {
					fmt.Println("Your loss limit has been reached.")
					fmt.Println("Stopping.")
					return balance
				}
				fmt.Println("Current balance is", balance)
				betTable = addNextFib(betTable[len(betTable)-2], betTable[len(betTable)-1], betTable)
				fmt.Println(betTable)
			}
		}
		betTable = []int{first, second}
	}

	return balance
}

func FibBettingSpec(baseBet, startBalance, endBalance, limitLoss float64, mode, high bool, apiKey, hash, curr string) float64 {
	const first int = 0
	const second int = 1
	const multiplier float64 = 1.2045
	var betTable []int = []int{first, second}
	balance := startBalance
	var betNumber uint = 1

	fmt.Println("Win multiplier is", multiplier)

	fmt.Println("Current balance is", balance)
	fmt.Println("Your base bet is", baseBet)
	if limitLoss == 0 {
		fmt.Println("Warning you can lose your entire balance!")
	}

	var choice string
	fmt.Print("Are you sure you wanna run fibonacci sequence <yes/no>? ")
	fmt.Scan(&choice)
	choice = strings.ToUpper(choice)
	if choice == "NO" || choice == "N" {
		return balance
	}

	for balance < endBalance && balance > limitLoss {
		fmt.Println(betTable)
		for len(betTable) > 1 {
			fmt.Printf("%d. bet:\n", betNumber)
			betNumber++
			currentBet := float64(betTable[len(betTable)-1]) * baseBet
			currentBet = math.Round(currentBet*100000) / 100000
			fmt.Println("Next bet is", currentBet)
			if currentBet > balance {
				fmt.Println("Balance too small!!!")
				return balance
			} else if balance-currentBet < limitLoss {
				fmt.Println("Your loss limit has been reached.")
				fmt.Println("Stopping.")
				return balance
			}

			// Placing a bet
			result, _ := PlaceABetSpec(apiKey, fmt.Sprintf("%.8f", currentBet), "44", curr, hash, mode, high, true)
			if result {
				fmt.Println("Success!✅")
				balance += currentBet * multiplier
				fmt.Println("Current balance is", balance)
				if balance >= endBalance {
					break
				}
				if len(betTable) == 2 {
					break
				} else {
					betTable = betTable[:len(betTable)-2]
					fmt.Println(betTable)
				}
			} else {
				fmt.Println("Failure!☯ ")
				balance -= currentBet
				if balance < limitLoss {
					fmt.Println("Your loss limit has been reached.")
					fmt.Println("Stopping.")
					return balance
				}
				fmt.Println("Current balance is", balance)
				betTable = addNextFib(betTable[len(betTable)-2], betTable[len(betTable)-1], betTable)
				fmt.Println(betTable)
			}
		}
		betTable = []int{first, second}
	}

	return balance
}
