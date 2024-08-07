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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

func main() {
	apiFile, err := os.Open("API")
	if err != nil {
		apiFile, err = os.Open("API.txt")
		if err != nil {
			log.Fatal(err)
		}
	}
	scanner := bufio.NewScanner(apiFile)
	scanner.Scan()
	apiKey := scanner.Text()
	url := "https://duckdice.io/api/bot/user-info?api_key=" + apiKey
	err = apiFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Create a CookieJar to store cookies
		jar, err := cookiejar.New(nil)
		if err != nil {
			fmt.Println("Error creating cookie jar:", err)
			log.Fatal(err)
		}

		// Create an HTTP client that uses the CookieJar
		client := &http.Client{
			Jar: jar,
		}
		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

		// Create a new HTTP GET request
		response, err := client.Get(url)
		if err != nil {
			fmt.Println("Error making request:", err)
			log.Fatal(err)
		}
		defer response.Body.Close() // Ensure the response body is closed later

		// Read the response data
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			log.Fatal(err)
		}

		// Retrying when captcha triggered
		for string(body)[2:9] == "DOCTYPE" {
			waiter := rand.Uint32N(3) + 1

			client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
			fmt.Println("CAPTCHA TRIGGERED!😠😠😠")
			fmt.Printf("Waiting %d seconds", waiter)

			// Implented up to 3 second wait
			for range waiter {
				time.Sleep(time.Second)
				fmt.Print(".")
			}
			fmt.Println()
			fmt.Println("Retrying...")

			//Send request again
			captchaResp, err := client.Get(url)
			if err != nil {
				fmt.Println("Error sending request:", err)
				log.Fatal(err)
			}
			defer captchaResp.Body.Close()

			// Read the new response
			body, err = io.ReadAll(captchaResp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				log.Fatal(err)
			}
		}

		var userInfo UserInfoResponse
		err = json.Unmarshal([]byte(body), &userInfo)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			log.Fatal(err)
		}

		var curr, balance, choice, bonusHash, tleHash string
		var baseBalance, bonusBalance float64
		faucet, isHigh, bonusM, bonusExist, tleExist := true, true, false, (len(userInfo.WageringBonuses) != 0), (len(userInfo.TLE) != 0)
		var tleM bool
		var strategy string

		fmt.Println("Username:", userInfo.Username)
		fmt.Println("-------------------------------")

		if tleExist {
			fmt.Println("Special events:")
			fmt.Println("-------------------------------")
			for _, tle := range userInfo.TLE {
				fmt.Println("   - Name:", tle.Name, " Status:", tle.Status, " Hash:", tle.Hash)
				tleHash = tle.Hash
			}
			fmt.Println("-------------------------------")
		}

		if bonusExist {
			fmt.Println("Wagering Bonuses:")
			fmt.Println("-------------------------------")
			for _, bonus := range userInfo.WageringBonuses {
				maxWin, _ := strconv.ParseFloat(bonus.Margin, 64)
				maxWin *= 3
				fmt.Println("   - Name:", bonus.Name, " Status:", bonus.Status, " Hash:", bonus.Hash,
					" Type:", bonus.Type, " Margin:", bonus.Margin, bonus.Symbol,
					" Max Win:", maxWin, bonus.Symbol)
				bonusHash = bonus.Hash
				curr = bonus.Symbol
				_, bonusBalance = PlaceABetSpec(apiKey, "0.0008", "95", curr, bonusHash, false, true, false)
			}
			fmt.Println("-------------------------------")
		}

		fmt.Println("Balances:")
		fmt.Println("-------------------------------")
		for _, balans := range userInfo.Balances {
			if balans.Main != "" {
				fmt.Println(balans.Main, " ", balans.Currency)
			}

			if balans.Faucet != "" {
				fmt.Println(balans.Faucet, " ", balans.Currency, "(Faucet)")
			}
		}
		if bonusExist {
			fmt.Println(bonusBalance, " ", curr, "(Bonus)")
		}
		fmt.Println("-------------------------------")

		fmt.Println("Betting strategies:")
		fmt.Println("-------------------------------")
		fmt.Println("[F]ibonacci")
		fmt.Println("[L]abouchere")
		fmt.Println("[O]ne Percent Hunt")
		fmt.Println("-------------------------------")

		fmt.Print("Which betting strategy would you like to use <fibonacci/labouchere/one percent hunt/exit>: ")
		fmt.Scan(&choice)

		switch choice {
		case "f", "F":
			strategy = "fibonacci"
		case "l", "L":
			strategy = "labouchere"
		case "o", "O":
			strategy = "onePercent"
		case "e", "E":
			fmt.Println("Goodbye.")
			os.Exit(0)
		default:
			fmt.Println("Goodbye.")
			os.Exit(0)
		}

		// Clearing curr so that bonus currency is not automatically selected.
		curr = ""

		var baseBet, targetBal, limitLoss float64

		fmt.Print("Which currency would you like to bet in: ")
		fmt.Scan(&curr)
		curr = strings.ToUpper(curr)
		fmt.Println("You chose", curr, "currency.")

		if curr == "DECOY" {
			faucet = false
			fmt.Println("Main mode automaticaly chosen for", curr, "betting.")
		} else {
			fmt.Print("Which mode would you like to bet in <faucet/main/bonus/tle>: ")
			fmt.Scan(&choice)
			if choice == "Main" || choice == "main" || choice == "M" || choice == "m" || choice == "MAIN" {
				faucet = false
				fmt.Println("You chose MAIN mode.")
			} else if choice == "Bonus" || choice == "bonus" || choice == "B" || choice == "b" || choice == "BONUS" {
				faucet = false
				bonusM = true
				fmt.Println("You chose BONUS mode.")
			} else if choice == "TLE" || choice == "Tle" || choice == "t" || choice == "T" || choice == "tle" {
				faucet = false
				tleM = true
				fmt.Println("You chose TLE mode.")
			}

		}

		if faucet {
			fmt.Println("You chose FAUCET mode.")
		}

		fmt.Print("Insert base bet value: ")
		fmt.Scan(&baseBet)
		if strategy == "labouchere" {
			fmt.Println("Max win:", baseBet*10, curr)
		}

		fmt.Print("Would you like to bet <high/low>: ")
		fmt.Scan(&choice)
		if choice == "Low" || choice == "low" || choice == "L" || choice == "l" || choice == "LOW" {
			isHigh = false
		}

		if !bonusM {
			for _, balans := range userInfo.Balances {
				if balans.Currency == curr {
					if faucet {
						balance = balans.Faucet
					} else {
						balance = balans.Main
					}
				}
			}
		}

		if bonusM {
			baseBalance = bonusBalance
		} else {
			baseBalance, _ = strconv.ParseFloat(balance, 64)
		}
		fmt.Printf("Balance is %.6f %s.\n", baseBalance, curr)

		if strategy == "labouchere" {
			fmt.Print("Insert target balance value: ")
			fmt.Scan(&targetBal)
			for targetBal-baseBalance > baseBet*10 {
				fmt.Println("Target balance too high. Look at max win above!!!")
				fmt.Print("Insert target balance value: ")
				fmt.Scan(&targetBal)
			}
			fmt.Printf("Target balance is %.6f %s.\n", targetBal, curr)
		} else if strategy == "fibonacci" {
			fmt.Print("Insert target balance value: ")
			fmt.Scan(&targetBal)
			fmt.Printf("Target balance is %.6f %s.\n", targetBal, curr)

			fmt.Print("Do you want to limit your loss <yes/no>? ")
			fmt.Scan(&choice)
			choice = strings.ToUpper(choice)
			if choice == "YES" || choice == "Y" {
				fmt.Print("Enter the minimum balance you would like to stop at: ")
				fmt.Scan(&limitLoss)
				fmt.Printf("In case of a bad streak the betting will stop when balance reaches %.6f %s.\n", limitLoss, curr)
			} else {
				limitLoss = 0
			}
		}

		if strategy == "labouchere" {
			if bonusM {
				temp := LabouchereSpec(baseBet, baseBalance, targetBal, faucet, isHigh, apiKey, bonusHash, curr)
				fmt.Println("Final balance is", temp, curr)
			} else if tleM {
				temp := LabouchereSpec(baseBet, baseBalance, targetBal, tleM, isHigh, apiKey, tleHash, curr)
				fmt.Println("Final balance is", temp, curr)
			} else {
				temp := Labouchere(baseBet, baseBalance, targetBal, faucet, isHigh, apiKey, curr)
				fmt.Println("Final balance is", temp, curr)
			}
		} else if strategy == "onePercent" {
			if bonusM {
				temp := OnePercentHuntSpec(baseBet, baseBalance, faucet, isHigh, apiKey, bonusHash, curr)
				fmt.Println("Final balance is", temp, curr)
			} else if tleM {
				temp := OnePercentHuntSpec(baseBet, baseBalance, tleM, isHigh, apiKey, tleHash, curr)
				fmt.Println("Final balance is", temp, curr)
			} else {
				temp := OnePercentHunt(baseBet, baseBalance, faucet, isHigh, apiKey, curr)
				fmt.Println("Final balance is", temp, curr)
			}
		} else if strategy == "fibonacci" {
			if bonusM {
				temp := FibBettingSpec(baseBet, baseBalance, targetBal, limitLoss, faucet, isHigh, apiKey, bonusHash, curr)
				fmt.Println("Final balance is", temp, curr)
			} else if tleM {
				temp := FibBettingSpec(baseBet, baseBalance, targetBal, limitLoss, tleM, isHigh, apiKey, tleHash, curr)
				fmt.Println("Final balance is", temp, curr)
			} else {
				temp := FibBetting(baseBet, baseBalance, targetBal, limitLoss, faucet, isHigh, apiKey, curr)
				fmt.Println("Final balance is", temp, curr)
			}

		}
	}
}
