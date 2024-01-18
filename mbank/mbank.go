package mbank

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/sanity-io/litter"
)

const mBankBaseUrl = "https://online.mbank.pl/pl"

type Client struct {
	countryCode string
	dfp         string
	mbank8      string
	userAgent   string
	jar         *cookiejar.Jar
}

func NewClient(mbank8, dfp, userAgent string) *Client {

	jr, _ := cookiejar.New(nil)

	return &Client{countryCode: "pl", mbank8: mbank8, dfp: dfp, userAgent: userAgent, jar: jr}
}

func (m *Client) ChangeProfile(company bool) (string, error) {

	payload := `{profileCode: "T"}`
	if company {
		payload = `{profileCode: "84491651"}`
	}

	operationsResponseRaw, err := m.request("/LoginMain/Account/JsonActivateProfile", strings.NewReader(payload), http.MethodPost)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Millisecond * 250)

	return operationsResponseRaw, nil
}

func (m *Client) AllAccounts() ([]Account, error) {

	tempMap := map[string]Account{}

	m.ChangeProfile(false)
	accounts, err := m.Accounts()
	if err != nil {
		return accounts, err
	}

	for _, a := range accounts {
		tempMap[string(a.AccountNumber)] = a
	}

	m.ChangeProfile(true)
	accounts2, _ := m.Accounts()
	if err != nil {
		return accounts, err
	}

	if len(accounts2) > 0 {
		if _, ok := tempMap[string(accounts2[0].AccountNumber)]; !ok {
			accounts = append(accounts, accounts2...)
		}
	}

	return accounts, nil
}

func (m *Client) Accounts() ([]Account, error) {
	var accounts []Account

	mainAccountsRawResponse, err := m.request("/Accounts/Accounts/AccountsGroups", nil, http.MethodGet)
	if err != nil {
		return accounts, err
	}

	savingAccountsRawResponse, err := m.request("/Accounts/Accounts/SavingAccounts", nil, http.MethodGet)
	if err != nil {
		return accounts, err
	}

	type AccountsGroupsRes struct {
		AccountsGroups []struct {
			Accounts []Account `json:"accounts"`
		} `json:"accountsGroups"`
	}

	parsedMainAccountResponse := &AccountsGroupsRes{}
	err = json.Unmarshal([]byte(mainAccountsRawResponse), parsedMainAccountResponse)

	if err != nil {
		return accounts, err
	}

	type SavingAccountsGroupsRes struct {
		Accounts []Account `json:"accounts"`
	}
	parsedSavingAccountResponse := &SavingAccountsGroupsRes{}

	err = json.Unmarshal([]byte(savingAccountsRawResponse), parsedSavingAccountResponse)
	if err != nil {
		return accounts, err
	}

	var accountsTemp []Account
	for _, groups := range parsedMainAccountResponse.AccountsGroups {
		accountsTemp = append(accountsTemp, groups.Accounts...)
	}

	accountsTemp = append(accountsTemp, parsedSavingAccountResponse.Accounts...)
	for _, a := range accountsTemp {
		accounts = append(accounts, Account{
			a.AccountNumber,
			a.Balance,
			a.Currency,
			a.Name,
			a.CustomName,
		})
	}

	return accounts, nil
}

func (m *Client) Login(login int, password string) error {

	dt := &LoginRequest{UserName: login, Password: password, Scenario: "Default", HrefHasHash: false}
	dt.DfpData = map[string]string{"dfp": m.dfp, "errorMessage": "", "scaOperationId": "5ff200dec4c44"}

	res, _ := json.Marshal(dt)
	body := string(res)

	resp, err := m.request("/Account/JsonLogin", strings.NewReader(body), http.MethodPost)
	if err != nil {
		return err
	}

	type LoginResponse struct {
		Successful bool `json:"successful"`
	}

	loginResponse := &LoginResponse{}
	err = json.Unmarshal([]byte(resp), loginResponse)
	if err != nil {
		return err
	}

	if !loginResponse.Successful {
		return errors.New("login failed")
	}

	return nil
}

func (m *Client) request(uri string, body io.Reader, method string) (string, error) {

	fullUrl := mBankBaseUrl + uri
	if uri == "/api/investmentfunds/v1/StockMarketAccount" {
		fullUrl = "https://online.mbank.pl/api/investmentfunds/v1/StockMarketAccount"
	}

	//fullUrl := "https://online.mbank.pl/api/investmentfunds/v1/StockMarketAccount"

	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return "", err
	}

	if false {
		debug := map[string]interface{}{
			"url":    fullUrl,
			"method": method,
		}

		litter.Dump(debug)
	}

	req.Header.Set("User-Agent", m.userAgent)
	if method == http.MethodPost {
		req.Header.Set("Cookie", "mBank8="+m.mbank8)
		req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	}

	if fullUrl == "https://online.mbank.pl/pl/Pfm/HistoryApi/GetHostTransactionsSummary" {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		req.Header.Set("Accept", "application/json, text/plain, */*")
	}

	client := http.Client{Jar: m.jar}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
