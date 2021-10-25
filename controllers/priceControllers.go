package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	easysdk "github.com/personal-security/easy-sdk-go"
	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

var CoinsPriceList = func(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	list, err := cg.CoinsList()
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	resp := easysdk.Message(true, "success")
	resp["items"] = list
	easysdk.Respond(w, resp)
}

var CoinsPriceListPrepared = func(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	list, err := cg.CoinsList()
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	items := make([]map[string]interface{}, 0)

	for _, item := range *list {
		tempItem := make(map[string]interface{})

		tempItem["$modelId"] = item.ID
		tempItem["$modelType"] = "Coin"
		tempItem["symbol"] = item.Symbol
		tempItem["name"] = item.ID

		items = append(items, tempItem)
	}

	resp := easysdk.Message(true, "Success")
	resp["items"] = items
	easysdk.Respond(w, resp)
}

var CoinsPriceListPreparedFiltered = func(w http.ResponseWriter, r *http.Request) {

	type ParamsPost struct {
		Coins []string
	}

	err := r.ParseForm()
	if err != nil {
		easysdk.Respond(w, easysdk.Message(false, err.Error()))
		return
	}

	decoder.IgnoreUnknownKeys(true)

	params := &ParamsPost{}

	// r.PostForm is a map of our POST form values
	err = decoder.Decode(params, r.PostForm)
	if err != nil {
		easysdk.Respond(w, easysdk.Message(false, err.Error()))
		return
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	list, err := cg.CoinsList()
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	ids := params.Coins
	vc := []string{"usd", "eur"}
	sp, err := cg.SimplePrice(ids, vc)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, "Timeout")
		easysdk.Respond(w, resp)
	}

	items := make([]map[string]interface{}, 0)

	for _, item := range *list {
		tempItem := make(map[string]interface{})

		tempItem["$modelId"] = item.ID
		tempItem["$modelType"] = "Coin"
		tempItem["symbol"] = item.Symbol
		tempItem["name"] = item.ID

		coinData := (*sp)[item.ID]

		ud := make(map[string]interface{})
		for _, v := range vc {
			ud[v] = coinData[v]
		}

		tempItem["prices"] = ud

		if easysdk.SliceHaveString(params.Coins, item.ID) {
			items = append(items, tempItem)
		}
	}

	resp := easysdk.Message(true, "Success")
	resp["items"] = items

	easysdk.Respond(w, resp)
}

func (client ClientHandler) CoinsPriceListPreparedFilteredJson(w http.ResponseWriter, r *http.Request) {

	type ParamsPost struct {
		Coins []string `json:"coins"`
	}

	params := &ParamsPost{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	var list *types.CoinList

	cachedList, isOk := client.Lru.Get("coins_list")
	if isOk {
		list = cachedList.(*types.CoinList)
	} else {
		list, err = cg.CoinsList()
		if err != nil {
			log.Println(err)
			resp := easysdk.Message(false, err.Error())
			easysdk.Respond(w, resp)
			return
		}
	}

	ids := params.Coins
	vc := []string{"usd", "eur"}
	sp, err := cg.SimplePrice(ids, vc)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, "Timeout")
		easysdk.Respond(w, resp)
	}

	items := make([]map[string]interface{}, 0)

	for _, item := range *list {
		tempItem := make(map[string]interface{})

		tempItem["$modelId"] = item.ID
		tempItem["$modelType"] = "Coin"
		tempItem["symbol"] = item.Symbol
		tempItem["name"] = item.ID

		coinData := (*sp)[item.ID]

		ud := make(map[string]interface{})
		for _, v := range vc {
			ud[v] = coinData[v]
		}

		tempItem["prices"] = ud

		if easysdk.SliceHaveString(params.Coins, item.ID) {
			items = append(items, tempItem)
		}
	}

	resp := easysdk.Message(true, "Success")
	resp["items"] = items

	easysdk.Respond(w, resp)
}

var CoinsPriceCurrent = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	coin := vars["coin"]

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	ids := []string{coin}
	vc := []string{"usd", "eur"}
	sp, err := cg.SimplePrice(ids, vc)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, "Timeout")
		easysdk.Respond(w, resp)
	}

	coinData := (*sp)[coin]
	//eth := (*sp)["ethereum"]

	ud := make(map[string]interface{})
	for _, v := range vc {
		ud[v] = coinData[v]
	}

	resp := easysdk.Message(true, "success")
	resp["item"] = ud
	easysdk.Respond(w, resp)
}

var CoinsPriceAll = func(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	ids := []string{"bitcoin", "ethereum"}
	vc := []string{"usd", "eur"}
	sp, err := cg.SimplePrice(ids, vc)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, "Timeout")
		easysdk.Respond(w, resp)
	}

	bitcoin := (*sp)["bitcoin"]
	eth := (*sp)["ethereum"]

	ud := make(map[string]interface{})
	ud["bitcoin"] = bitcoin["usd"]
	ud["ethereum"] = eth["usd"]

	resp := easysdk.Message(true, "success")
	resp["items"] = ud
	easysdk.Respond(w, resp)
}
