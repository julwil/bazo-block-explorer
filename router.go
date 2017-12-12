package main

import (
    "io"
    "fmt"
    "os/exec"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/dgrijalva/jwt-go"
    "github.com/mchetelat/bazo_miner/protocol"
)

func initializeRouter() *httprouter.Router {
  router := httprouter.New()
  router.GET("/", getIndex)
  //router.GET("/get-token", getToken
  //router.GET("/testheader", testSPVHeader)
  router.GET("/testblock", testBlock)
  router.GET("/blocks", getAllBlocks)
  router.GET("/block/:hash", getOneBlock)
  router.GET("/transactions/funds", getAllFundsTx)
  router.GET("/transactions/funds/:hash", getOneFundsTx)
  router.GET("/transactions/acc", getAllAccTx)
  router.GET("/transactions/acc/:hash", getOneAccTx)
  router.GET("/transactions/config", getAllConfigTx)
  router.GET("/transactions/config/:hash", getOneConfigTx)
  router.GET("/account/:hash", getAccount)
  router.POST("/search/", searchForHash)
  router.POST("/login", loginFunc)
  router.GET("/login-failed", loginFail)
  router.GET("/adminpanel", adminfunc)
  router.POST("/admin/blocksize", blocksizePost)
  router.POST("/admin/diffinterval", diffintervalPost)
  router.POST("/admin/minfee", minfeePost)
  router.POST("/admin/blockinterval", blockintervalPost)
  router.POST("/admin/blockreward", blockrewardPost)
  router.ServeFiles("/static/*filepath", http.Dir("static"))

  return router
}

func testBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  //genesishash := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }
  var block *protocol.Block = reqBlock(nil)
  tpl.ExecuteTemplate(w, "realblock.gohtml", block)
}
/*
func testSPVHeader(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  var header *protocol.SPVHeader = reqSPVHeader(nil)
  fmt.Printf("%x\n", header.Hash)
  tpl.ExecuteTemplate(w, "header.gohtml", header)
}
*/
func getToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  cookie := CreateToken()
  http.SetCookie(w, &cookie)
  http.Redirect(w, r, "/", 307)
}

func getIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedrows := ReturnBlocksAndTransactions(params)
  tpl.ExecuteTemplate(w, "index.gohtml", returnedrows)
}

func getOneBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblock := ReturnOneBlock(params)
  tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)

}
func getAllBlocks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblocks := ReturnAllBlocks(params)
  tpl.ExecuteTemplate(w, "blocks.gohtml", returnedblocks)
}

func getOneFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := ReturnOneFundsTx(params)
  /*
  fmt.Println(returnedtx.Hash)
  if returnedtx.Hash == "" {
    fmt.Println("trying to copy opentx from network")
    txHash := params.ByName("hash")
    FetchOpenTx(txHash)
    returnedtx = ReturnOpenFundsTx(params)
  }
  */
  tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedtx)
}

func getAllFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := ReturnAllFundsTx(params)
  tpl.ExecuteTemplate(w, "fundstxs.gohtml", returnedtxs)
}

func getOneAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := ReturnOneAccTx(params)
  tpl.ExecuteTemplate(w, "acctx.gohtml", returnedtx)
}

func getAllAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := ReturnAllAccTx(params)
  tpl.ExecuteTemplate(w, "acctxs.gohtml", returnedtxs)
}

func getOneConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := ReturnOneConfigTx(params)
  tpl.ExecuteTemplate(w, "configtx.gohtml", returnedtx)
}

func getAllConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := ReturnAllConfigTx(params)
  tpl.ExecuteTemplate(w, "configtxs.gohtml", returnedtxs)
}

func searchForHash(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  thing := 1
  returnedblock, returnedtx := ReturnSearchResult(r)
  if returnedblock.Hash == "" && returnedtx.Hash == "" {
    tpl.ExecuteTemplate(w, "noresult.gohtml", thing)
  } else if returnedblock.Hash != "" && returnedtx.Hash == "" {
    tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)
  } else if returnedblock.Hash == "" && returnedtx.Hash != "" {
    tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedtx)
  } else {
    tpl.ExecuteTemplate(w, "noresult.gohtml", thing)
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

func blocksizeGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "blocksize.gohtml", values)
}

func blocksizePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  payload := r.PostFormValue("new-blocksize")
  cmd := exec.Command("bazo_client", "configTx", "0", "1", payload, "1", "2", "root")
  printCommand(cmd)
  output, err := cmd.CombinedOutput()
  printError(err)
  printOutput(output)
  io.WriteString(w, string(output))
}

func diffintervalGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "difficulty.gohtml", values)
}

func diffintervalPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  payload := r.PostFormValue("new-diffinterval")
  cmd := exec.Command("bazo_client", "configTx", "0", "2", payload, "1", "2", "root")
  printCommand(cmd)
  output, err := cmd.CombinedOutput()
  printError(err)
  printOutput(output)
  io.WriteString(w, string(output))
}

func minfeeGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "minfee.gohtml", values)
}

func minfeePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  payload := r.PostFormValue("new-minfee")
  cmd := exec.Command("bazo_client", "configTx", "0", "3", payload, "2", "4", "root")
  printCommand(cmd)
  output, err := cmd.CombinedOutput()
  printError(err)
  printOutput(output)
  io.WriteString(w, string(output))
}

func blockintervalGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "blockinterval.gohtml", values)
}

func blockintervalPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  payload := r.PostFormValue("new-blockinterval")
  cmd := exec.Command("bazo_client", "configTx", "0", "4", payload, "1", "2", "root")
  printCommand(cmd)
  output, err := cmd.CombinedOutput()
  printError(err)
  printOutput(output)
  io.WriteString(w, string(output))
}

func blockrewardGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "blockreward.gohtml", values)
}

func blockrewardPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  payload := r.PostFormValue("new-blockreward")
  cmd := exec.Command("bazo_client", "configTx", "0", "5", payload, "1", "2", "root")
  printCommand(cmd)
  output, err := cmd.CombinedOutput()
  printError(err)
  printOutput(output)
  io.WriteString(w, string(output))
}

func getTestPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  thing := 1
  tpl.ExecuteTemplate(w, "test.gohtml", thing)
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
  returnedaccount := ReturnOneAccount(params)
  tpl.ExecuteTemplate(w, "account.gohtml", returnedaccount)
}
