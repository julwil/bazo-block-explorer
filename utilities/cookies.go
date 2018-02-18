package utilities

import (
  "net/http"
  "time"
  "fmt"
  "io/ioutil"
  "encoding/json"
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
  RESTresponse, err := http.Get("http://192.41.136.199:443/account/" + publicKey)
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
