package controllers

import (
	"coin-nodes-rest/models"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	easysdk "github.com/personal-security/easy-sdk-go"
)

var CoinsList = func(w http.ResponseWriter, r *http.Request) {
	items := make([]string, 0)

	items = append(items, "eth")
	items = append(items, "btc")

	resp := easysdk.Message(true, "success")
	resp["items"] = items
	easysdk.Respond(w, resp)
}

var NodesInfoAll = func(w http.ResponseWriter, r *http.Request) {

}

var NodeInfo = func(w http.ResponseWriter, r *http.Request) {

}

var WalletGenerate = func(w http.ResponseWriter, r *http.Request) {

}

func (client ClientHandler) WalletInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	module := vars["type_wallet"]
	address := vars["wallet"]

	switch module {
	case "eth":
		if address == "" {
			resp := easysdk.Message(false, "Malformed request")
			easysdk.Respond(w, resp)
			return
		}

		balance, err := models.GetAddressBalance(*client.Eth, address)

		if err != nil {
			fmt.Println(err)
			resp := easysdk.Message(false, "Internal server error")
			easysdk.Respond(w, resp)
			return
		}

		balanceFloat, err := easysdk.StringToFloat64(balance)
		if err != nil {
			resp := easysdk.Message(false, "Bad balance")
			easysdk.Respond(w, resp)
			return
		}

		walletInfo := &models.WalletInfo{
			Address: address,
			Amount:  balanceFloat,
		}

		resp := easysdk.Message(true, "Success")
		resp["item"] = walletInfo
		easysdk.Respond(w, resp)
		return
	case "btc":
		balanceFloat := models.BtcWalletGetBalance(address)
		balanceUnconfirmedFloat := models.BtcWalletGetBalanceUnconfirmed(address)

		walletInfo := &models.WalletInfo{
			Address:           address,
			Amount:            balanceFloat,
			AmountUnconfirmed: balanceUnconfirmedFloat,
		}

		resp := easysdk.Message(true, "Success")
		resp["item"] = walletInfo
		easysdk.Respond(w, resp)
		return
	}

	resp := easysdk.Message(false, "Internal server error")
	easysdk.Respond(w, resp)
}

func (client ClientHandler) TransactionInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	module := vars["type_wallet"]
	tx := vars["tx"]

	switch module {
	case "eth":
		transactionInfo := &models.TransactionInfo{
			Tx: tx,
		}

		resp := easysdk.Message(true, "Success")
		resp["item"] = transactionInfo
		easysdk.Respond(w, resp)
		return
	case "btc":

	}

	resp := easysdk.Message(false, "Internal server error")
	easysdk.Respond(w, resp)
}

var MoneySend = func(w http.ResponseWriter, r *http.Request) {

	type ParamsPost struct {
		To     string  `json:"to"`
		From   string  `json:"from"`
		Amount float64 `json:"amount"`
		Fee    float64 `json:"fee"`
	}

	vars := mux.Vars(r)

	module := vars["type_wallet"]

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

	//sms.Text = r.FormValue("text")

	switch module {
	case "eth":

	case "btc":
		if models.BtcSendMoney(params.From, params.To, params.Amount) {
			resp := easysdk.Message(true, "Success")
			easysdk.Respond(w, resp)
			return
		} else {
			resp := easysdk.Message(false, "Error send bitcoin")
			easysdk.Respond(w, resp)
		}
	}

	resp := easysdk.Message(false, "Internal server error")
	easysdk.Respond(w, resp)
}

var MoneyMultiSend = func(w http.ResponseWriter, r *http.Request) {

}

var ToolsAllNodesRestart = func(w http.ResponseWriter, r *http.Request) {

}

var ToolsNodeRestart = func(w http.ResponseWriter, r *http.Request) {

}

func (client ClientHandler) BroadcastRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	module := vars["type_wallet"]

	type ParamsPost struct {
		RawTx string `json:"raw_tx"`
	}

	params := &ParamsPost{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Println(err)
		resp := easysdk.Message(false, err.Error())
		easysdk.Respond(w, resp)
		return
	}

	params.RawTx = strings.TrimPrefix(params.RawTx, "0x")

	switch module {
	case "eth":
		rawTxBytes, err := hex.DecodeString(params.RawTx)
		if err != nil {
			log.Println(err)
			resp := easysdk.Message(false, err.Error())
			easysdk.Respond(w, resp)
			return
		}

		tx := new(types.Transaction)
		rlp.DecodeBytes(rawTxBytes, &tx)

		err = client.Eth.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Println(err)
			resp := easysdk.Message(false, err.Error())
			easysdk.Respond(w, resp)
			return
		}

		resp := easysdk.Message(true, "Success")
		easysdk.Respond(w, resp)
		return
	case "btc":

	}

	resp := easysdk.Message(false, "Internal server error")
	easysdk.Respond(w, resp)
}

func (client ClientHandler) TransactionHelper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	module := vars["type_wallet"]

	switch module {
	case "eth":
		type SendData struct {
			Nonce    uint64  `json:"nonce"`
			GasLimit int     `json:"gasLimit"`
			GasPrice big.Int `json:"gasPrice"`
			ChainID  big.Int `json:"chainId"`
		}

		type EthQuery struct {
			From string `json:"from"`
			To   string `json:"to"`
		}

		params := &EthQuery{}

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			log.Println(err)
			resp := easysdk.Message(false, err.Error())
			easysdk.Respond(w, resp)
			return
		}

		var gasLimit int

		tokenAddressTo := common.HexToAddress(params.To)
		estimatedGas, err := client.Eth.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &tokenAddressTo,
			Data: []byte{0},
		})
		if err != nil {
			log.Println(err)
			gasLimit = 0
		}

		accelerationVariable := 1.0 //1.30
		gasLimit = int(float64(estimatedGas) * accelerationVariable)

		tokenAddressFrom := common.HexToAddress(params.From)
		nonce, err := client.Eth.PendingNonceAt(context.Background(), tokenAddressFrom)
		if err != nil {
			log.Fatal(err)
		}
		gasPrice, err := client.Eth.SuggestGasPrice(context.Background())
		if err != nil {
			resp := easysdk.Message(false, err.Error())
			easysdk.Respond(w, resp)
			return
		}

		chainID, err := client.Eth.NetworkID(context.Background())
		if err != nil {
			resp := easysdk.Message(false, err.Error())
			easysdk.Respond(w, resp)
			return
		}

		item := &SendData{
			Nonce:    nonce,
			GasLimit: gasLimit,
			GasPrice: *gasPrice,
			ChainID:  *chainID,
		}

		resp := easysdk.Message(true, "Success")
		resp["item"] = item
		easysdk.Respond(w, resp)
		return
	case "btc":
		type SendData struct {
			Fee float64 `json:"fee"`
		}

		item := &SendData{
			Fee: models.BtcWalletGetFee(),
		}

		resp := easysdk.Message(true, "Success")
		resp["item"] = item
		easysdk.Respond(w, resp)
		return
	}

	resp := easysdk.Message(false, "Internal server error")
	easysdk.Respond(w, resp)
}
