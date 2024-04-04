package main

import "math/big"

type LastDeposit struct {
	CreatedAt int    `json:"createdAt"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
}

type CurrencyAmount struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type Balance struct {
	Currency string `json:"currency"`
	Main     string `json:"main"`
	Faucet   string `json:"faucet"`
}

type WageringBonus struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Hash   string `json:"hash"`
	Status string `json:"status"`
	Symbol string `json:"symbol"`
	Margin string `json:"margin"`
}

type TLE struct {
	Hash   string `json:"hash"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UserInfoResponse struct {
	Hash            string           `json:"hash"`
	Username        string           `json:"username"`
	CreatedAt       int              `json:"createdAt"`
	Campaign        string           `json:"campaign"`
	Affiliate       string           `json:"affiliate"`
	Level           int              `json:"level"`
	LastDeposit     LastDeposit      `json:"lastDeposit"`
	Wagered         []CurrencyAmount `json:"wagered"`
	Balances        []Balance        `json:"balances"`
	WageringBonuses []WageringBonus  `json:"wageringBonuses"`
	TLE             []TLE            `json:"tle"`
}

type Bet struct {
	Hash      string  `json:"hash"`
	Symbol    string  `json:"symbol"`
	Result    bool    `json:"result"`
	IsHigh    bool    `json:"isHigh"`
	Number    int     `json:"number"`
	Treshold  int     `json:"treshold"`
	Chance    float64 `json:"chance"`
	Payout    float64 `json:"payout"`
	BetAmount string  `json:"betAmount"`
	WinAmount string  `json:"winAmount"`
	Profit    string  `json:"profit"`
	Mined     string  `json:"mined"`
	Nonce     int     `json:"nonce"`
	Created   int     `json:"created"`
	GameMode  string  `json:"gameMode"`
}

type AbsoluteLevel struct {
	Level  int `json:"level"`
	Xp     int `json:"xp"`
	XpNext int `json:"xpNext"`
	XpPrev int `json:"xpPrev"`
}

type User struct {
	Hash          string        `json:"hash"`
	Level         int           `json:"level"`
	Username      string        `json:"username"`
	Bets          int           `json:"bets"`
	Nonce         int           `json:"nonce"`
	Wins          int           `json:"wins"`
	Luck          float64       `json:"luck"`
	Balance       *big.Float    `json:"balance"` // Using big.Float for precision
	Profit        *big.Float    `json:"profit"`  // Using big.Float for precision
	Volume        *big.Float    `json:"volume"`  // Using big.Float for precision
	AbsoluteLevel AbsoluteLevel `json:"absoluteLevel"`
}

type Jackpot struct {
	Amount string `json:"amount"`
	User   User   `json:"user"`
}

type BetResponse struct {
	Bet           Bet      `json:"bet"`
	IsJackpot     bool     `json:"isJackpot"`
	JackpotStatus *bool    `json:"jackpotStatus"`
	Jackpot       *Jackpot `json:"jackpot"`
	User          User     `json:"user"`
}

// BetPayload struct to represent the JSON data
type BetPayload struct {
	Symbol string `json:"symbol"`
	Chance string `json:"chance"`
	IsHigh bool   `json:"isHigh"`
	Amount string `json:"amount"`
	Faucet bool   `json:"faucet"`
}
