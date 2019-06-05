package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var tokenName string

func CreateCookie(publicKey string) http.Cookie {
	tokenName = "publicKey"
	expireCookie := time.Now().Add(time.Minute * 10)

	cookie := http.Cookie{Name: tokenName, Value: publicKey, Expires: expireCookie}
	return cookie
}

func GetPublicKeyCookie(r *http.Request) (*http.Cookie, error) {
	tokenName := "publicKey"
	publicKeyCookie, err := r.Cookie(tokenName)
	return publicKeyCookie, err
}

func RequestAccountInformation(publicKey string) JSONAccount {
	//IP address is client
	RESTresponse, err := http.Get("https://oysyone.westeurope.cloudapp.azure.com/account/" + publicKey)
	if err != nil {
		fmt.Print(err.Error())
	}

	var accountInformation JSONAccount
	RESTresponseData, err := ioutil.ReadAll(RESTresponse.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(RESTresponseData, &accountInformation)

	return accountInformation
}
