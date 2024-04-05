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
		return
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
		return
	}
	defer response.Body.Close() // Ensure the response body is closed later

	// Inspect received cookies
	//	cookies := response.Cookies()
	//	fmt.Println("Received Cookies:", cookies)

	// Read the response data
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(27) + 3

		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
		fmt.Println("FUCKING CAPTCHA!😠😠😠")
		fmt.Printf("Waiting %d seconds", waiter)

		// Implented 10 second wait
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
			return
		}
		defer captchaResp.Body.Close()

		// Read the new response
		body, err = io.ReadAll(captchaResp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
	}

	var userInfo UserInfoResponse
	err = json.Unmarshal([]byte(body), &userInfo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Username:", userInfo.Username)
	fmt.Println("Balances:")
	for _, balans := range userInfo.Balances {
		if balans.Main == "" {
			fmt.Println(balans.Faucet, " ", balans.Currency, "(Faucet)")
		} else {
			fmt.Println(balans.Main, " ", balans.Currency)
		}
	}

	var bet float64
	var currency string
	var faucetMode string = "faucet"
	fMode := true
	fmt.Print("Insert bet value: ")
	fmt.Scan(&bet)
	amount := fmt.Sprintf("%f", bet)
	fmt.Println("Bet is", amount)

	fmt.Print("Choose currency: ")
	fmt.Scan(&currency)

	fmt.Print("Choose mode: ")
	fmt.Scan(&faucetMode)

	if faucetMode == "Main" || faucetMode == "main" {
		fMode = false
	}

	rez := PlaceABet(apiKey, amount, currency, fMode)
	if rez == true {
		fmt.Println("Bet successful.✅")
	} else {
		fmt.Println("Bet unsuccessful.☯")
	}
}
