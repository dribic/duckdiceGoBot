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
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(apiFile)
	scanner.Scan()
	apiKey := scanner.Text() // Replace with your actual API key
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

	// Inspect received cookies
	//	cookies := response.Cookies()
	//	fmt.Println("Received Cookies:", cookies)

	// Read the response data
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		log.Fatal(err)
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(6)

		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
		fmt.Println("FUCKING CAPTCHA!😠😠😠")
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

	//test stuff
	var curr, balance string

	fmt.Println("Username:", userInfo.Username)
	fmt.Println("Balances:")
	fmt.Println("-------------------------------")
	for _, balans := range userInfo.Balances {
		if balans.Main == "" {
			fmt.Println(balans.Faucet, " ", balans.Currency, "(Faucet)")
		} else {
			fmt.Println(balans.Main, " ", balans.Currency)
		}
		//test stuff
		if balans.Currency == "USDT" {
			curr = balans.Currency
			balance = balans.Faucet
		}
	}
	fmt.Println("-------------------------------")

	var baseBet, targetBal float64
	fmt.Print("Insert base bet value: ")
	fmt.Scan(&baseBet)
	fmt.Println("Max win:", baseBet*10, curr)

	//test stuff
	balD, _ := strconv.ParseFloat(balance, 64)
	fmt.Printf("Balance is %.6f %s. balD var type is %T!\n", balD, curr, balD)

	fmt.Print("Insert target balance value: ")
	fmt.Scan(&targetBal)
	for targetBal-balD > baseBet*10 {
		fmt.Println("Target balance too high. Look at max win above!!!")
		fmt.Print("Insert target balance value: ")
		fmt.Scan(&targetBal)
	}
	fmt.Printf("Target balance is %.6f %s.\n", targetBal, curr)

	Labouchere(baseBet, balD, targetBal, true, true, apiKey, curr)
}
