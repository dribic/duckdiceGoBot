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

import "fmt"

func OnePercentHunt(startBet, startBalance float64, fMode, high bool, apiKey, curr string) float64 {
	var currentBalance float64 = startBalance
	currentBet, totalLoss := startBet, startBet
	var choice string

	for range 98 {
		currentBet *= 1.01
		totalLoss += currentBet
	}

	// Resetting currentBet
	currentBet = startBet

	fmt.Println("If you start betting with", startBet, curr, "you can lose", totalLoss, curr)
	fmt.Print("Are you sure you want to continue? <yes/no> ")
	fmt.Scan(&choice)
	if choice == "No" || choice == "NO" || choice == "n" || choice == "N" || choice == "no" {
		return currentBalance
	}

	for i := range 99 {
		fmt.Printf("%d. bet:\n", i+1)
		if currentBet > currentBalance {
			fmt.Println("Balance is too low!")
			break
		}
		currentBalance -= currentBet
		fmt.Println("Current bet is", currentBet)
		fmt.Println("Current balance is", currentBalance)

		// Turn bet into a string, because DDs API uses string
		amount := fmt.Sprintf("%.8f", currentBet)

		// Making a bet
		result := PlaceABet(apiKey, amount, "0.95", curr, fMode, high)

		if result {
			if fMode {
				currentBalance += currentBet * 102.11
			} else {
				currentBalance += currentBet * 104.21
			}
			fmt.Println("Success!✅")
			return currentBalance
		} else {
			fmt.Println("Failure!☯")
		}

		currentBet *= 1.01
	}

	return currentBalance
}

func OnePercentHuntSpec(startBet, startBalance float64, mode, high bool, apiKey, hash, curr string) float64 {
	var currentBalance float64 = startBalance
	currentBet, totalLoss := startBet, startBet
	var choice string

	for range 98 {
		currentBet *= 1.01
		totalLoss += currentBet
	}

	// Resetting currentBet
	currentBet = startBet

	fmt.Println("If you start betting with", startBet, curr, "you can lose", totalLoss, curr)
	fmt.Print("Are you sure you want to continue? <yes/no> ")
	fmt.Scan(&choice)
	if choice == "No" || choice == "NO" || choice == "n" || choice == "N" || choice == "no" {
		return currentBalance
	}

	for i := range 99 {
		fmt.Printf("%d. bet:\n", i+1)
		if currentBet > currentBalance {
			fmt.Println("Balance is too low!")
			break
		}
		currentBalance -= currentBet
		fmt.Println("Current bet is", currentBet)
		fmt.Println("Current balance is", currentBalance)

		// Turn bet into a string, because DDs API uses string
		amount := fmt.Sprintf("%.8f", currentBet)

		// Making a bet
		result, _ := PlaceABetSpec(apiKey, amount, "0.95", curr, hash, mode, high, true)

		if result {
			currentBalance += currentBet * 102.11
			fmt.Println("Success!✅")
			return currentBalance
		} else {
			fmt.Println("Failure!☯")
		}

		currentBet *= 1.01
	}

	return currentBalance
}
