package main

import (
    "fmt"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/dgrijalva/jwt-go"
)

func initializeRouter() *httprouter.Router {
  router := httprouter.New()

  router.GET("/", getIndex)
  //router.GET("/get-token", getToken
  router.GET("/blocks", getAllBlocks)
  router.GET("/block/:hash", getOneBlock)
  router.GET("/transactions/funds", getAllFundsTx)
  router.GET("/transactions/funds/:hash", getOneFundsTx)
  router.GET("/transactions/acc", getAllAccTx)
  router.GET("/transactions/acc/:hash", getOneAccTx)
  router.GET("/transactions/config", getAllConfigTx)
  router.GET("/transactions/config/:hash", getOneConfigTx)
  router.GET("/account/:hash", getAccount)
  router.GET("/accounts", getTopAccounts)
  router.POST("/search/", searchForHash)
  router.POST("/login", loginFunc)
  router.GET("/login-failed", loginFail)
  router.GET("/adminpanel", adminNoVerif)
  router.ServeFiles("/static/*filepath", http.Dir("static"))

  return router
}

func adminNoVerif(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  tpl.ExecuteTemplate(w, "admin.gohtml", values)
}

func getToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  cookie := CreateToken()
  http.SetCookie(w, &cookie)
  http.Redirect(w, r, "/", 307)
}

func getIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedrows := ReturnBlocksAndTransactions(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "index.gohtml", returnedrows)
}

func getOneBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblock := ReturnOneBlock(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)

}
func getAllBlocks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblocks := ReturnAllBlocks(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "blocks.gohtml", returnedblocks)
}

func getOneFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := ReturnOneFundsTx(params.ByName("hash"))
  /*
  fmt.Println(returnedtx.Hash)
  if returnedtx.Hash == "" {
    fmt.Println("trying to copy opentx from network")
    txHash := params.ByName("hash")
    FetchOpenTx(txHash)
    returnedtx = ReturnOpenFundsTx(params.ByName("hash"))
  }
  */
  tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedtx)
}

func getAllFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := ReturnAllFundsTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "fundstxs.gohtml", returnedtxs)
}

func getOneAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := ReturnOneAccTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "acctx.gohtml", returnedtx)
}

func getAllAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := ReturnAllAccTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "acctxs.gohtml", returnedtxs)
}

func getOneConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := ReturnOneConfigTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "configtx.gohtml", returnedtx)
}

func getAllConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := ReturnAllConfigTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "configtxs.gohtml", returnedtxs)
}

func searchForHash(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblock := ReturnOneBlock(r.PostFormValue("search-value"))
  if returnedblock.Hash != "" {
    tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)
  }

  returnedfundstx := ReturnOneFundsTx(r.PostFormValue("search-value"))
  if returnedfundstx.Hash != "" {
    tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedfundstx)
  }

  returnedacctx := ReturnOneAccTx(r.PostFormValue("search-value"))
  if returnedacctx.Hash != "" {
    tpl.ExecuteTemplate(w, "acctx.gohtml", returnedacctx)
  }

  returnedconfigtx := ReturnOneConfigTx(r.PostFormValue("search-value"))
  if returnedconfigtx.Hash != "" {
    tpl.ExecuteTemplate(w, "configtx.gohtml", returnedconfigtx)
  }

  returnedaccountwithtxs := ReturnOneAccount(r.PostFormValue("search-value"))
  if returnedaccountwithtxs.Account.Hash != "" {
    tpl.ExecuteTemplate(w, "account.gohtml", returnedaccountwithtxs)
  }
}

func adminfunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tokenCookie, err := ExtractCookie(r)
	switch {
	case err == http.ErrNoCookie:
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "No cookie in request!")
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while Parsing cookie!")
		fmt.Fprintln(w, "Cookie parse error: %v\n")
		return
	}
  token, err := ParseToken(tokenCookie)
  switch err.(type) {
	case nil: // no error
		if !token.Valid { // but may still be invalid
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid Token")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
    tpl.ExecuteTemplate(w, "admin.gohtml", values)

	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Token Expired, get a new one.")
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while Parsing Token!")
			fmt.Printf("ValidationError error: %+v\n", vErr.Errors)
			return
		}

	default: // something else went wrong
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while Parsing Token!")
		fmt.Printf("Token parse error: %v\n", err)
		return
	}
}

func loginFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  if r.PostFormValue("root-key-field") == "123456" {
    cookie := CreateToken()
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/", 302)
    //http.Redirect(w, r, "/get-token", 302)
  } else {
    http.Redirect(w, r, "/login-failed", 302)
  }
}

func loginFail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  thing := 1
  tpl.ExecuteTemplate(w, "loginfail.gohtml", thing)
}

func getAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccountwithtxs := ReturnOneAccount(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "account.gohtml", returnedaccountwithtxs)
}

func getTopAccounts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccounts := ReturnTopAccounts(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "accounts.gohtml", returnedaccounts)
}
