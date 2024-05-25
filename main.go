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
		waiter := rand.Uint32N(5) + 1

		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
		fmt.Println("CAPTCHA TRIGGERED!😠😠😠")
		fmt.Printf("Waiting %d seconds", waiter)

		// Implented up to 15 second wait
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

	var curr, balance, choice, hash string
	var baseBalance, bonusBalance float64
	faucet, isHigh, bonusM, bonusExist, tleExist := true, true, false, (len(userInfo.WageringBonuses) != 0), (len(userInfo.TLE) != 0)

	fmt.Println("Username:", userInfo.Username)
	fmt.Println("-------------------------------")

	if tleExist {
		fmt.Println("Special events:")
		fmt.Println("-------------------------------")
		for _, tle := range userInfo.TLE {
			fmt.Println("   - Name:", tle.Name, " Status:", tle.Status, " Hash:", tle.Hash)
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
			hash = bonus.Hash
			curr = bonus.Symbol
			_, bonusBalance = PlaceABetSpec(apiKey, "0.0008", "95", curr, hash, false, true, false)
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

	// Clearing curr so that bonus currency is not automatically selected.
	curr = ""

	var baseBet, targetBal float64
	var progSteps uint8 = 1
	var progress string = "no"

	fmt.Print("Which currency would you like to bet in: ")
	fmt.Scan(&curr)
	fmt.Println("You chose", curr, "currency.")

	fmt.Print("Which mode would you like to bet in <faucet/main/bonus>: ")
	fmt.Scan(&choice)
	if choice == "Main" || choice == "main" || choice == "M" || choice == "m" || choice == "MAIN" {
		faucet = false
	} else if choice == "Bonus" || choice == "bonus" || choice == "B" || choice == "b" || choice == "BONUS" {
		faucet = false
		bonusM = true
	}

	fmt.Print("Insert base bet value: ")
	fmt.Scan(&baseBet)
	fmt.Println("Max win:", baseBet*10, curr)

	fmt.Print("Do you want progressive betting <yes/no>? ")
	fmt.Scan(&progress)

	if progress == "Yes" || progress == "yes" || progress == "Y" || progress == "y" || progress == "YES" {
		fmt.Print("How many steps do you want? ")
		fmt.Scan(&progSteps)
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

	fmt.Print("Insert target balance value: ")
	fmt.Scan(&targetBal)
	for targetBal-baseBalance > baseBet*10 {
		fmt.Println("Target balance too high. Look at max win above!!!")
		fmt.Print("Insert target balance value: ")
		fmt.Scan(&targetBal)
	}
	fmt.Printf("Target balance is %.6f %s.\n", targetBal, curr)

	if bonusM {
		if progSteps == 1 {
			temp := LabouchereSpec(baseBet, baseBalance, targetBal, faucet, isHigh, apiKey, hash, curr)
			fmt.Println("Final balance is", temp, curr)
		} else {
			temp := baseBalance
			for i := range progSteps {
				if i > 0 {
					fmt.Println("-------------------------------")
				}
				fmt.Printf("%d. step:\n", i+1)
				fmt.Println("-------------------------------")
				baseBalance = LabouchereSpec(baseBet, temp, (targetBal + baseBet*10*float64(i)), faucet, isHigh, apiKey, hash, curr)
				temp = baseBalance
				fmt.Println("Success!✅")
				if i == 0 {
					fmt.Println("-------------------------------")
				}
			}
			fmt.Println("Final balance is", temp, curr)
		}
	} else {
		if progSteps == 1 {
			temp := Labouchere(baseBet, baseBalance, targetBal, faucet, isHigh, apiKey, curr)
			fmt.Println("Final balance is", temp, curr)
		} else {
			temp := baseBalance
			for i := range progSteps {
				if i > 0 {
					fmt.Println("-------------------------------")
				}
				fmt.Printf("%d. step:\n", i+1)
				fmt.Println("-------------------------------")
				baseBalance = Labouchere(baseBet, temp, (targetBal + baseBet*10*float64(i)), faucet, isHigh, apiKey, curr)
				temp = baseBalance
				fmt.Println("Success!✅")
				if i == 0 {
					fmt.Println("-------------------------------")
				}
			}
			fmt.Println("Final balance is", temp, curr)
		}
	}
}
