package router

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/way365/bazo-block-explorer/data"
	"github.com/way365/bazo-block-explorer/utilities"
	"html/template"
	"net/http"
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
	router.GET("/tx/update", getAllUpdateTx)
	router.GET("/tx/update/:hash", getOneUpdateTx)
	router.GET("/tx/agg", getAllAggTx)
	router.GET("/tx/agg/:hash", getOneAggTx)
	router.GET("/tx/config", getAllConfigTx)
	router.GET("/tx/config/:hash", getOneConfigTx)
	router.GET("/tx/stake", getAllStakeTx)
	router.GET("/tx/stake/:hash", getOneStakeTx)
	router.GET("/account/:hash", getAccount)
	router.GET("/accounts", getTopAccounts)
	router.GET("/stats", getStats)
	router.GET("/search", searchForHash)
	router.POST("/login", loginFunc)
	router.GET("/logout", logoutFunc)
	router.GET("/adminpanel", adminfunc)

	router.ServeFiles("/source/*filepath", http.Dir("./source"))

	return router
}

func getIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedrows := data.ReturnBlocksAndTransactions(params.ByName("hash"))
	returnedrows.UrlLevel = ""
	tpl.ExecuteTemplate(w, "index.gohtml", returnedrows)
}

func getOneBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedblock := data.ReturnOneBlock(params.ByName("hash"))
	returnedblock.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "block.gohtml", returnedblock)

}
func getAllBlocks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var blocks utilities.Blocksandurl
	blocks.Blocks = data.ReturnAllBlocks(params.ByName("hash"))
	blocks.UrlLevel = ".."
	tpl.ExecuteTemplate(w, "blocks.gohtml", blocks)
}

func getOneFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedtx := data.ReturnOneFundsTx(params.ByName("hash"))
	returnedtx.UrlLevel = "../../.."
	tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedtx)
}

func getAllFundsTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var txs utilities.Fundsandurl
	txs.Txs = data.ReturnAllFundsTx(params.ByName("hash"))
	txs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "fundstxs.gohtml", txs)
}

func getOneAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedtx := data.ReturnOneAccTx(params.ByName("hash"))
	returnedtx.UrlLevel = "../../.."
	tpl.ExecuteTemplate(w, "acctx.gohtml", returnedtx)
}

func getAllUpdateTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var txs utilities.Updatesandurl
	txs.Txs = data.ReturnAllUpdateTx(params.ByName("hash"))
	txs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "updatetxs.gohtml", txs)
}

func getOneUpdateTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedtx := data.ReturnOneUpdateTx(params.ByName("hash"))
	returnedtx.UrlLevel = "../../.."
	tpl.ExecuteTemplate(w, "updatetx.gohtml", returnedtx)
}

func getAllAggTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var txs utilities.AggsAndUrl
	txs.Txs = data.ReturnAllAggTx(params.ByName("hash"))
	txs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "aggtxs.gohtml", txs)
}

func getOneAggTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedtx := data.ReturnOneAggTx(params.ByName("hash"))
	returnedtx.UrlLevel = "../../.."
	tpl.ExecuteTemplate(w, "aggtx.gohtml", returnedtx)
}

func getAllAccTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var txs utilities.Accsandurl
	txs.Txs = data.ReturnAllAccTx(params.ByName("hash"))
	txs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "acctxs.gohtml", txs)
}

func getOneConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedtx := data.ReturnOneConfigTx(params.ByName("hash"))
	returnedtx.UrlLevel = "../../.."
	tpl.ExecuteTemplate(w, "configtx.gohtml", returnedtx)
}

func getAllConfigTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var txs utilities.Configsandurl
	txs.Txs = data.ReturnAllConfigTx(params.ByName("hash"))
	txs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "configtxs.gohtml", txs)
}

func getOneStakeTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedtx := data.ReturnOneStakeTx(params.ByName("hash"))
	returnedtx.UrlLevel = "../../.."
	tpl.ExecuteTemplate(w, "staketx.gohtml", returnedtx)
}

func getAllStakeTx(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var txs utilities.Stakesandurl
	txs.Txs = data.ReturnAllStakeTx(params.ByName("hash"))
	txs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "staketxs.gohtml", txs)
}

func searchForHash(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	queryValues := r.URL.Query()
	hash := queryValues.Get("hash")

	returnedaccountwithtxs := data.ReturnOneAccount(hash)
	fmt.Printf("Query Parameter: %v", hash)
	if returnedaccountwithtxs.Account.Hash != "" {
		returnedaccountwithtxs.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "accountSearch.gohtml", returnedaccountwithtxs)
		return
	}

	returnedconfigtx := data.ReturnOneConfigTx(hash)
	if returnedconfigtx.Hash != "" {
		returnedconfigtx.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "configtx.gohtml", returnedconfigtx)
		return
	}

	returnedblock := data.ReturnOneBlock(hash)
	if returnedblock.Hash != "" {
		returnedblock.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "blockSearch.gohtml", returnedblock)
		return
	}

	returnedfundstx := data.ReturnOneFundsTx(hash)
	if returnedfundstx.Hash != "" {
		returnedfundstx.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "fundstx.gohtml", returnedfundstx)
		return
	}

	returnedacctx := data.ReturnOneAccTx(hash)
	if returnedacctx.Hash != "" {
		returnedacctx.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "acctx.gohtml", returnedacctx)
		return
	}

	returnedstaketx := data.ReturnOneStakeTx(hash)
	if returnedstaketx.Hash != "" {
		returnedstaketx.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "staketx.gohtml", returnedstaketx)
		return
	}
	tpl.ExecuteTemplate(w, "noresult.gohtml", returnedstaketx)
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
		parameters.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "admin.gohtml", parameters)
	} else {
		var url utilities.Emptyresponse
		url.UrlLevel = ".."
		tpl.ExecuteTemplate(w, "loginfail.gohtml", url)
	}

}

func loginFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	accountInformation := utilities.RequestAccountInformation(r.PostFormValue("public-key-field"))
	var url utilities.Emptyresponse
	url.UrlLevel = ".."

	if accountInformation.IsRoot {
		cookie := utilities.CreateCookie(r.PostFormValue("public-key-field"))
		http.SetCookie(w, &cookie)
		tpl.ExecuteTemplate(w, "loginsuccess.gohtml", url)
	} else {
		tpl.ExecuteTemplate(w, "loginfail.gohtml", url)
	}
}

func logoutFunc(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	cookie := utilities.CreateCookie(" ")
	var url utilities.Emptyresponse
	url.UrlLevel = ".."

	http.SetCookie(w, &cookie)
	tpl.ExecuteTemplate(w, "loggedout.gohtml", url)
}

func getAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnedaccountwithtxs := data.ReturnOneAccount(params.ByName("hash"))
	for _, tx := range returnedaccountwithtxs.Txs {
		tx.UrlLevel = "../.."
	}
	returnedaccountwithtxs.UrlLevel = "../.."
	tpl.ExecuteTemplate(w, "account.gohtml", returnedaccountwithtxs)
}

func getTopAccounts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var accounts utilities.Accountsandurl
	accounts.Accounts = data.ReturnTopAccounts(params.ByName("hash"))
	accounts.UrlLevel = ".."
	tpl.ExecuteTemplate(w, "accounts.gohtml", accounts)
}

func getStats(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	stats := data.ReturnTotals()
	stats.Parameters = data.ReturnNewestParameters()
	chartData := data.Return14Hours()
	b, _ := json.Marshal(chartData)
	stats.ChartData = string(b)
	stats.UrlLevel = ".."

	tpl.ExecuteTemplate(w, "stats.gohtml", stats)
}
