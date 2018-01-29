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
  router.GET("/blocks", getAllBlocks)
  router.GET("/block/:hash", getOneBlock)
  router.GET("/tx/funds", getAllFundsTx)
  router.GET("/tx/funds/:hash", getOneFundsTx)
  router.GET("/tx/acc", getAllAccTx)
  router.GET("/tx/acc/:hash", getOneAccTx)
  router.GET("/tx/config", getAllConfigTx)
  router.GET("/tx/config/:hash", getOneConfigTx)
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

func getIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedrows := data.ReturnBlocksAndTransactions(params.ByName("hash"))
  returnedrows.UrlLevel = ""
  tpl.ExecuteTemplate(w, "index.gohtml", returnedrows)
}

func getOneBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblock := data.ReturnOneBlock(params.ByName("hash"))
  returnedblock.UrlLevel = "../../"
  tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)

}
func getAllBlocks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedblocks := data.ReturnAllBlocks(params.ByName("hash"))
  for _, block := range returnedblocks {
    block.UrlLevel = "../"
  }
  tpl.ExecuteTemplate(w, "blocks.gohtml", returnedblocks)
}

func getOneFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := data.ReturnOneFundsTx(params.ByName("hash"))
  returnedtx.UrlLevel = "../../../"
  tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedtx)
}

func getAllFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := data.ReturnAllFundsTx(params.ByName("hash"))
  for _, tx := range returnedtxs {
    tx.UrlLevel = "../../"
  }
  tpl.ExecuteTemplate(w, "fundstxs.gohtml", returnedtxs)
}

func getOneAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := data.ReturnOneAccTx(params.ByName("hash"))
  returnedtx.UrlLevel = "../../../"
  tpl.ExecuteTemplate(w, "acctx.gohtml", returnedtx)
}

func getAllAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := data.ReturnAllAccTx(params.ByName("hash"))
  for _, tx := range returnedtxs {
    tx.UrlLevel = "../../"
  }
  tpl.ExecuteTemplate(w, "acctxs.gohtml", returnedtxs)
}

func getOneConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtx := data.ReturnOneConfigTx(params.ByName("hash"))
  returnedtx.UrlLevel = "../../../"
  tpl.ExecuteTemplate(w, "configtx.gohtml", returnedtx)
}

func getAllConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedtxs := data.ReturnAllConfigTx(params.ByName("hash"))
  for _, tx := range returnedtxs {
    tx.UrlLevel = "../../"
  }
  tpl.ExecuteTemplate(w, "configtxs.gohtml", returnedtxs)
}

func searchForHash(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccountwithtxs := data.ReturnOneAccount(r.PostFormValue("search-value"))
  if returnedaccountwithtxs.Account.Hash != "" {
    returnedaccountwithtxs.UrlLevel = "../"
    tpl.ExecuteTemplate(w, "accountSearch.gohtml", returnedaccountwithtxs)
    return
  }

  returnedconfigtx := data.ReturnOneConfigTx(r.PostFormValue("search-value"))
  if returnedconfigtx.Hash != "" {
    returnedconfigtx.UrlLevel = "../"
    tpl.ExecuteTemplate(w, "configtxSearch.gohtml", returnedconfigtx)
    return
  }

  returnedblock := data.ReturnOneBlock(r.PostFormValue("search-value"))
  if returnedblock.Hash != "" {
    returnedblock.UrlLevel = "../"
    tpl.ExecuteTemplate(w, "blockSearch.gohtml", returnedblock)
    return
  }

  returnedfundstx := data.ReturnOneFundsTx(r.PostFormValue("search-value"))
  if returnedfundstx.Hash != "" {
    returnedfundstx.UrlLevel = "../"
    tpl.ExecuteTemplate(w, "fundstxSearch.gohtml", returnedfundstx)
    return
  }

  returnedacctx := data.ReturnOneAccTx(r.PostFormValue("search-value"))
  if returnedacctx.Hash != "" {
    returnedacctx.UrlLevel = "../"
    tpl.ExecuteTemplate(w, "acctxSearch.gohtml", returnedacctx)
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
    parameters.UrlLevel = "../"
    tpl.ExecuteTemplate(w, "admin.gohtml", parameters)
  } else {
    urlLevel := "../"
    tpl.ExecuteTemplate(w, "loginfail.gohtml", urlLevel)
  }

}

func loginFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  accountInformation := utilities.RequestAccountInformation(r.PostFormValue("public-key-field"))
  urlLevel := "../"

  if accountInformation.IsRoot {
    cookie := utilities.CreateCookie(r.PostFormValue("public-key-field"))
    http.SetCookie(w, &cookie)
    tpl.ExecuteTemplate(w, "loginsuccess.gohtml", urlLevel)
  } else {
    tpl.ExecuteTemplate(w, "loginfail.gohtml", urlLevel)
  }
}

func logoutFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  cookie := utilities.CreateCookie(" ")
  urlLevel := "../"
  http.SetCookie(w, &cookie)
  tpl.ExecuteTemplate(w, "loggedout.gohtml", urlLevel)
}

func getAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccountwithtxs := data.ReturnOneAccount(params.ByName("hash"))
  for _, tx := range returnedaccountwithtxs.Txs {
    tx.UrlLevel = "../../"
  }
  returnedaccountwithtxs.UrlLevel = "../../"
  tpl.ExecuteTemplate(w, "account.gohtml", returnedaccountwithtxs)
}

func getTopAccounts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccounts := data.ReturnTopAccounts(params.ByName("hash"))
  for _, account := range returnedaccounts {
    account.UrlLevel = "../"
  }
  tpl.ExecuteTemplate(w, "accounts.gohtml", returnedaccounts)
}

func getStats(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  stats := data.ReturnTotals()
  stats.Parameters = data.ReturnNewestParameters()
  chartData := data.Return14Hours()
  b, _ := json.Marshal(chartData)
  stats.ChartData = string(b)
  stats.UrlLevel = "../"

  tpl.ExecuteTemplate(w, "stats.gohtml", stats)
}
