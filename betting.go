package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"time"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

func PlaceABet(apiKey, betValue, currency string) bool {
	var result bool
	url := "https://duckdice.io/api/play?api_key=" + apiKey

	// Create a bet
	sampleBet := BetPayload{
		Symbol: currency,
		Chance: "44",
		IsHigh: false,
		Amount: betValue,
		Faucet: true,
	}

	// Create a CookieJar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return false
	}

	// Marshal the payload into JSON
	jsonPayload, err := json.Marshal(sampleBet)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return false
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
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
		return false
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return false
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(27) + 3
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return false
		}

		// Set the Content-Type header
		req.Header.Set("Content-Type", "application/json")

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
		captchaResp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return false
		}
		defer captchaResp.Body.Close()

		// Read the new response
		body, err = io.ReadAll(captchaResp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return false
		}
	}

	var betResp BetResponse
	err = json.Unmarshal([]byte(body), &betResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false
	}

	result = betResp.Bet.Result

	fmt.Println("Roll is", betResp.Bet.Number)

	return result
}
