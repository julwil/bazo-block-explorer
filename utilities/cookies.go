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
	//IP address is bootstrap server
	RESTresponse, err := http.Get("http://localhost:8010/account/" + publicKey)
	if err != nil {
		fmt.Print(err.Error())
	}

	responseRaw, err := ioutil.ReadAll(RESTresponse.Body)
	if err != nil {
		panic(err)
	}

	var responseBody JSONAccountResponseBody
	var account JSONAccount

	json.Unmarshal(responseRaw, &responseBody)

	if len(responseBody.Content) == 0 {
		return account
	}

	return responseBody.Content[0].Detail
}
