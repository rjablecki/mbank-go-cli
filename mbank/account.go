package mbank

import "strings"

type Currency string

const PLN Currency = "PLN"
const EUR Currency = "EUR"

type Account struct {
	AccountNumber AccountNumber `json:"accountNumber"`
	Balance       float32       `json:"balance"`
	Currency      Currency      `json:"currency"`
	Name          string        `json:"name"`
	CustomName    string        `json:"customName"`
}

func (a *Account) GetNumberShort() string {
	return strings.ReplaceAll(a.AccountNumber.String(), " ", "")
}

type AccountNumber string

func (b *AccountNumber) String() string {
	return string(*b)
}
