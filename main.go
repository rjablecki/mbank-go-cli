package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"rjablecki/mbank-cli-go/mbank"

	"github.com/joho/godotenv"
	"github.com/sanity-io/litter"
)

func getEnv(key string) string {

	return os.Getenv(key)
}

func main() {
	godotenv.Load(".env.local", ".env")
	mbankClient, err := mBankLogin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	balance, err := mbankClient.AllAccounts()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonResponse, _ := json.Marshal(balance)
	fmt.Println(string(jsonResponse))
}

func mBankLogin() (*mbank.Client, error) {

	mBankCredentialsBase := getEnv("MBANK_CREDENTIALS")
	mBankCredentials, err := base64.StdEncoding.DecodeString(mBankCredentialsBase)
	if err != nil {
		litter.Dump(err)
		return nil, err
	}
	var conf mbank.Config
	json.Unmarshal(mBankCredentials, &conf)

	client := mbank.NewClient(conf.Mbank8, conf.Dfp, conf.UserAgent)
	err = client.Login(conf.Id, conf.Pass)
	if err != nil {
		litter.Dump(err)
		return nil, err
	}
	return client, nil
}
