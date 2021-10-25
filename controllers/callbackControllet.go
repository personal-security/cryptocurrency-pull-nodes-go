package controllers

import (
	"coin-nodes-rest/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	easysdk "github.com/personal-security/easy-sdk-go"
)

var CallbackGet = func(w http.ResponseWriter, r *http.Request) {

	type QueryGet struct {
		NetWork string `json:"network"`
	}

	queryGet := &QueryGet{}

	err := json.NewDecoder(r.Body).Decode(&queryGet)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	vars := mux.Vars(r)

	module := vars["type_wallet"]
	wallet := vars["wallet"]

	callbackWallet := models.CallbackWalletGet(module, queryGet.NetWork, wallet)

	if callbackWallet == nil {
		easysdk.GenerateApiError(w, "Record not found", nil, 404)
	}

	easysdk.GenerateApiRespond(w, true, "Success", structs.Map(callbackWallet))

}

var CallbackSet = func(w http.ResponseWriter, r *http.Request) {

	type QueryUpdate struct {
		Url     string `json:"url"`
		NetWork string `json:"network"`
	}

	queryUpdate := &QueryUpdate{}

	err := json.NewDecoder(r.Body).Decode(&queryUpdate)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	vars := mux.Vars(r)

	module := vars["type_wallet"]
	wallet := vars["wallet"]

	callbackWallet := &models.CallbackWallet{}
	callbackWallet.Node = module
	callbackWallet.Network = queryUpdate.NetWork
	callbackWallet.Address = wallet
	callbackWallet.Url = queryUpdate.Url

	status, message := callbackWallet.Create()
	easysdk.GenerateApiRespond(w, status, message, structs.Map(callbackWallet))
}
