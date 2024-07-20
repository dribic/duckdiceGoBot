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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

func PlaceABet(apiKey, betValue, chance, currency string, mode, high bool) bool {
	var result bool
	url := "https://duckdice.io/api/play?api_key=" + apiKey

	// Create a bet
	bet := BetPayload{
		Symbol: currency,
		Chance: chance,
		IsHigh: high,
		Amount: betValue,
		Faucet: mode,
	}

	// Create a CookieJar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		os.Exit(1)
	}

	// Marshal the payload into JSON
	jsonPayload, err := json.Marshal(bet)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		os.Exit(1)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client that uses the CookieJar
	client := &http.Client{
		Jar: jar,
	}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(3) + 1
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("Error creating request:", err)
			os.Exit(1)
		}

		// Set the Content-Type header
		req.Header.Set("Content-Type", "application/json")

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
		captchaResp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			os.Exit(1)
		}
		defer captchaResp.Body.Close()

		// Read the new response
		body, err = io.ReadAll(captchaResp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}
	}

	var betResp BetResponse
	err = json.Unmarshal([]byte(body), &betResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	result = betResp.Bet.Result

	fmt.Println("Roll is", betResp.Bet.Number)
	fmt.Println("Payout is", betResp.Bet.WinAmount)

	return result
}
