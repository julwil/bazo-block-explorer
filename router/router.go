package router

import (
  "github.com/lucBoillat/BazoBlockExplorer/data"
  "github.com/lucBoillat/BazoBlockExplorer/utilities"
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/julienschmidt/httprouter"
  "html/template"
)

var tpl *template.Template

func InitializeRouter() *httprouter.Router {
  tpl = template.Must(template.ParseGlob("source/html/*"))

  router := httprouter.New()

  router.GET("/", getIndex)
  router.GET("/test", getTest)
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
  router.GET("/stats", getStats)
  router.POST("/search/", searchForHash)
  router.POST("/login", loginFunc)
  router.GET("/logout", logoutFunc)
  router.GET("/adminpanel", adminfunc)

  router.ServeFiles("/source/*filepath", http.Dir("source"))

  return router
}

func getTest(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  fmt.Println("testconsole")
  fmt.Fprintln(w, "testpage")
  return
}

func getIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedrows := data.ReturnBlocksAndTransactions(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "index.gohtml", returnedrows)
}

func getOneBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblock := data.ReturnOneBlock(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)

}
func getAllBlocks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblocks := data.ReturnAllBlocks(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "blocks.gohtml", returnedblocks)
}

func getOneFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := data.ReturnOneFundsTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedtx)
}

func getAllFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := data.ReturnAllFundsTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "fundstxs.gohtml", returnedtxs)
}

func getOneAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := data.ReturnOneAccTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "acctx.gohtml", returnedtx)
}

func getAllAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := data.ReturnAllAccTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "acctxs.gohtml", returnedtxs)
}

func getOneConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := data.ReturnOneConfigTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "configtx.gohtml", returnedtx)
}

func getAllConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := data.ReturnAllConfigTx(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "configtxs.gohtml", returnedtxs)
}

func searchForHash(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccountwithtxs := data.ReturnOneAccount(r.PostFormValue("search-value"))
  if returnedaccountwithtxs.Account.Hash != "" {
    tpl.ExecuteTemplate(w, "account.gohtml", returnedaccountwithtxs)
    return
  }

  returnedconfigtx := data.ReturnOneConfigTx(r.PostFormValue("search-value"))
  if returnedconfigtx.Hash != "" {
    tpl.ExecuteTemplate(w, "configtx.gohtml", returnedconfigtx)
    return
  }

  returnedblock := data.ReturnOneBlock(r.PostFormValue("search-value"))
  if returnedblock.Hash != "" {
    tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)
    return
  }

  returnedfundstx := data.ReturnOneFundsTx(r.PostFormValue("search-value"))
  if returnedfundstx.Hash != "" {
    tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedfundstx)
    return
  }

  returnedacctx := data.ReturnOneAccTx(r.PostFormValue("search-value"))
  if returnedacctx.Hash != "" {
    tpl.ExecuteTemplate(w, "acctx.gohtml", returnedacctx)
    return
  }
  tpl.ExecuteTemplate(w, "noresult.gohtml", returnedacctx)
}

func adminfunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  publicKeyCookie, err := utilities.GetPublicKeyCookie(r)
	switch {
  case err == http.ErrNoCookie:
    w.WriteHeader(http.StatusUnauthorized)
    fmt.Fprintln(w, "No cookie in request!")
    return
  case publicKeyCookie.Value == " ":
    w.WriteHeader(http.StatusUnauthorized)
    fmt.Fprintln(w, "Not verified!")
    return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while Parsing cookie!")
		fmt.Fprintln(w, "Cookie parse error: %v\n")
		return
	}

  accountInformation := utilities.RequestAccountInformation(publicKeyCookie.Value)
  if accountInformation.IsRoot {
    parameters := data.ReturnNewestParameters()
    tpl.ExecuteTemplate(w, "admin.gohtml", parameters)
  } else {
    tpl.ExecuteTemplate(w, "loginfail.gohtml", 1)
  }

}

func loginFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  accountInformation := utilities.RequestAccountInformation(r.PostFormValue("public-key-field"))

  if accountInformation.IsRoot {
    cookie := utilities.CreateCookie(r.PostFormValue("public-key-field"))
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/adminpanel", 302)
  } else {
    tpl.ExecuteTemplate(w, "loginfail.gohtml", 1)
  }
}

func logoutFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  cookie := utilities.CreateCookie(" ")
  http.SetCookie(w, &cookie)
  http.Redirect(w, r, "/", 302)
}

func getAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccountwithtxs := data.ReturnOneAccount(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "account.gohtml", returnedaccountwithtxs)
}

func getTopAccounts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccounts := data.ReturnTopAccounts(params.ByName("hash"))
  tpl.ExecuteTemplate(w, "accounts.gohtml", returnedaccounts)
}

func getStats(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  stats := data.ReturnTotals()
  stats.Parameters = data.ReturnNewestParameters()
  chartData := data.Return14Hours()
  b, _ := json.Marshal(chartData)
  stats.ChartData = string(b)

  tpl.ExecuteTemplate(w, "stats.gohtml", stats)
}
