package main

import (
    "net/http"
    "time"
    "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secret")
var tokenName = "AuthToken"

func CreateToken() http.Cookie {
  mySigningKey := []byte("secret")
  tokenName := "AuthToken"

  expireToken := time.Now().Add(time.Minute * 10).Unix()
  expireCookie := time.Now().Add(time.Minute * 10)
  token := jwt.New(jwt.SigningMethodHS256)
  Claims := token.Claims.(jwt.MapClaims)
  Claims["admin"] = true
  Claims["name"] = "root"
  Claims["exp"] = expireToken
  tokenString, _ := token.SignedString(mySigningKey)
  cookie := http.Cookie{Name: tokenName, Value: tokenString, Expires: expireCookie, HttpOnly: true}
  return cookie
}

func ExtractCookie(r *http.Request) (*http.Cookie, error) {
  tokenCookie, err := r.Cookie(tokenName)
  return tokenCookie, err
}

func ParseToken(tokenCookie *http.Cookie) (*jwt.Token, error) {
  token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
  return mySigningKey, nil
  })
  return token, err
}
