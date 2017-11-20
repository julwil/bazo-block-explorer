package main

import (
  "io"
  _ "fmt"
  "net/http"
  "html/template"
  _ "database/sql"
  _ "github.com/lib/pq"
  "github.com/julienschmidt/httprouter"
  _ "strconv"
  _ "strings"
)

var tpl *template.Template

var values systemparams

func init() {
  tpl = template.Must(template.ParseGlob("static/src/*.gohtml"))
  values.BlockSize = 1
  values.BSName = "Block Size"
  values.DiffInterval = 1
  values.DIName = "Difficulty Interval"
  values.MinFee = 1
  values.MFName = "Minimum Fee"
  values.BlockInterval = 1
  values.BIName = "Block Interval"
  values.BlockReward = 1
  values.BRName = "Block Reward"
}

func main() {
  router := httprouter.New()
  router.GET("/", getIndex)
  router.GET("/test", getTestPage)
  router.GET("/blocks", getAllBlocks)
  router.GET("/block/:hash", getOneBlock)
  router.GET("/transactions/funds", getAllFundsTx)
  router.GET("/transactions/fundtx/:hash", getOneFundsTx)
  router.GET("/transactions/acc", getAllAccTx)
  router.GET("/transactions/acc/:hash", getOneAccTx)
  router.GET("/transactions/config", getAllConfigTx)
  router.GET("/transactions/config/:hash", getOneConfigTx)
  router.GET("/account/:hash", getAccount)
  router.POST("/search/", searchForHash)
  router.GET("/adminpanel", adminfunc)
  router.GET("/admin/blocksize", blocksizeGet)
  router.POST("/admin/blocksize", blocksizePost)
  router.GET("/admin/diffinterval", diffintervalGet)
  router.POST("/admin/diffinterval", diffintervalPost)
  router.GET("/admin/minfee", minfeeGet)
  router.POST("/admin/minfee", minfeePost)
  router.GET("/admin/blockinterval", blockintervalGet)
  router.POST("/admin/blockinterval", blockintervalPost)
  router.GET("/admin/blockreward", blockrewardGet)
  router.POST("/admin/blockreward", blockrewardPost)
  router.ServeFiles("/static/*filepath", http.Dir("static"))
  http.ListenAndServe(":8080", router)
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

func getAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  returnedaccount := ReturnAccount(params)
  tpl.ExecuteTemplate(w, "account.gohtml", returnedaccount)
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
  tpl.ExecuteTemplate(w, "admin.gohtml", values)
}

func blocksizeGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "blocksize.gohtml", values)
}

func blocksizePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  textinput := r.PostFormValue("new-blocksize")
  io.WriteString(w, textinput)
}

func diffintervalGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "difficulty.gohtml", values)
}

func diffintervalPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  textinput := r.PostFormValue("new-diffinterval")
  io.WriteString(w, textinput)
}

func minfeeGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "minfee.gohtml", values)
}

func minfeePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  textinput := r.PostFormValue("new-minfee")
  io.WriteString(w, textinput)
}

func blockintervalGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "blockinterval.gohtml", values)
}

func blockintervalPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  textinput := r.PostFormValue("new-blockinterval")
  io.WriteString(w, textinput)
}

func blockrewardGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  tpl.ExecuteTemplate(w, "blockreward.gohtml", values)
}

func blockrewardPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    panic(err)
  }
  textinput := r.PostFormValue("new-blockreward")
  io.WriteString(w, textinput)
}

func getTestPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  thing := 1
  tpl.ExecuteTemplate(w, "test.gohtml", thing)
}
