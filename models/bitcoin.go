package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

type BitcoinConfigDB struct {
	BitcoinRPCUser string `json:"bitcoin_rpc_user"`
	BitcoinRPCPass string `json:"bitcoin_rpc_pass"`
	BitcoinRPCPort string `json:"bitcoin_rpc_port"`
}

type BtcAddressesUnspend struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
	Txid    string  `json:"txid"`
	Vout    int     `json:"vout"`
}

type BtcServerRequest struct {
	Method string        `json:"method"`
	Id     int           `json:"id"`
	Params []interface{} `json:"params"`
}

func BtcInit() *BitcoinConfigDB {
	config := &BitcoinConfigDB{}
	file, _ := ioutil.ReadFile("config/bitcoin.json")
	json.Unmarshal([]byte(file), &config)
	return config
}

func BtcRoundAmount(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		//rounder = math.Ceil(intermed)
		rounder = math.Floor(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

func (tp *BtcAddressesUnspend) UnmarshalJSON(data []byte) error {

	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		fmt.Printf("Error whilde decoding %v\n", err)
		return err
	}

	for k, v := range jsonMap {
		switch k {
		case "address":
			tp.Address = v.(string)
		case "amount":
			tp.Amount = v.(float64)
		case "txid":
			tp.Txid = v.(string)
		case "vout":
			tp.Vout = int(v.(float64))
		default:

		}
	}

	return nil
}

func BtcWalletGetFee() float64 {
	return 0.0002
}

func BtcRequestApi(Query string) []byte {
	config := BtcInit()
	url := "http://" + config.BitcoinRPCUser + ":" + config.BitcoinRPCPass + "@127.0.0.1:" + config.BitcoinRPCPort

	resp, err := http.Post(url,
		"application/json", strings.NewReader(Query))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
		return nil
	}

	return body
}

func BtcListUnspend(min uint, max uint, addresses []string) []BtcAddressesUnspend {
	type ServerAnswerListUnspend struct {
		Error  string                `json:"error"`
		Id     uint                  `json:"id"`
		Result []BtcAddressesUnspend `json:"result"`
	}

	request := BtcServerRequest{}
	request.Method = "listunspent"
	request.Id = 1
	request.Params = []interface{}{min, max, addresses}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswerListUnspend{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return make([]BtcAddressesUnspend, 0)
	}

	return result.Result
}

func BtcWalletValudate(address string) bool {

	type RequestValid struct {
		IsValid bool `json:"isvalid"`
	}

	type ServerAnswer struct {
		Error  string       `json:"error"`
		Id     uint         `json:"id"`
		Result RequestValid `json:"result"`
	}

	request := BtcServerRequest{}
	request.Method = "validateaddress"
	request.Id = 1
	request.Params = []interface{}{address}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswer{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return false
	}

	return result.Result.IsValid
}

func BtcWalletGetFreeBalance(address string) float64 {
	return BtcWalletGetBalance(address) - BtcWalletGetBalanceUnconfirmed(address)
}

func BtcWalletGetBalance(address string) float64 {
	var maxBalance float64 = 0

	data := BtcListUnspend(0, 99999, []string{address})
	for _, value := range data {
		maxBalance += value.Amount
	}

	return maxBalance
}

func BtcWalletGetBalanceUnconfirmed(address string) float64 {
	var maxBalance float64 = 0

	data := BtcListUnspend(0, 0, []string{address})
	for _, value := range data {
		maxBalance += value.Amount
	}

	return maxBalance
}

func BtcCreateRawTransaction(listAddresses []BtcAddressesUnspend, fromWallet string, toWallet string, amount float64) string {

	type ServerAnswerCreateRawTransaction struct {
		Error  string `json:"error"`
		Id     uint   `json:"id"`
		Result string `json:"result"`
	}

	type SendTo struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}

	type AddressesSendFrom struct {
		Txid string `json:"txid"`
		Vout int    `json:"vout"`
	}

	var maxBalance float64 = 0

	countUT := 0
	for _, value := range listAddresses {
		if maxBalance > amount {
			//continue
		}
		maxBalance = maxBalance + value.Amount
		countUT++

	}

	sendFrom := make([]AddressesSendFrom, countUT)
	for i, value := range listAddresses {
		if i >= countUT {
			//continue
		}
		sF := AddressesSendFrom{}
		sF.Txid = value.Txid
		sF.Vout = value.Vout
		sendFrom[i] = sF

	}

	toSend := []map[string]float64{}
	toSend = append(toSend, map[string]float64{toWallet: BtcRoundAmount(amount-BtcWalletGetFee(), 8)})

	if (maxBalance - amount) > BtcWalletGetFee() {
		toSend = append(toSend, map[string]float64{fromWallet: BtcRoundAmount(maxBalance-amount, 8)})
	}

	request := BtcServerRequest{}
	request.Method = "createrawtransaction"
	request.Id = 1
	request.Params = []interface{}{sendFrom, toSend}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswerCreateRawTransaction{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		//log.Fatalf("Unmarshal: %v", err)
		log.Printf("Unmarshal: %v", err)
		return ""
	}

	return result.Result
}

func BtcCreateRawTransactionMixerRand(listAddresses []BtcAddressesUnspend, fromWallet string, toWallet []string, amount float64) string {

	type ServerAnswerCreateRawTransaction struct {
		Error  string `json:"error"`
		Id     uint   `json:"id"`
		Result string `json:"result"`
	}

	type SendTo struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}

	type AddressesSendFrom struct {
		Txid string `json:"txid"`
		Vout int    `json:"vout"`
	}

	var maxBalance float64 = 0

	countUT := 0
	for _, value := range listAddresses {
		if maxBalance > amount {
			//continue
		}
		maxBalance = maxBalance + value.Amount
		countUT++

	}

	sendFrom := make([]AddressesSendFrom, countUT)
	for i, value := range listAddresses {
		if i >= countUT {
			//continue
		}
		sF := AddressesSendFrom{}
		sF.Txid = value.Txid
		sF.Vout = value.Vout
		sendFrom[i] = sF

	}

	toSend := []map[string]float64{}
	//amount_to_send := float64(amount)
	new_balande := amount
	amount_to_send := float64(0)
	for i := len(toWallet); i > 0; i-- {
		if i > 1 {
			amount_to_send = new_balande - (new_balande / float64(i) / 2)
		} else {
			amount_to_send = new_balande
		}
		if i == len(toWallet) {
			toSend = append(toSend, map[string]float64{toWallet[i-1]: BtcRoundAmount(amount_to_send-BtcWalletGetFee(), 8)})
		} else {
			toSend = append(toSend, map[string]float64{toWallet[i-1]: BtcRoundAmount(amount_to_send, 8)})
		}
		new_balande = new_balande - amount_to_send
		log.Printf("Balance %d send %f", i, new_balande)
	}

	if (maxBalance - amount) > BtcWalletGetFee() {
		toSend = append(toSend, map[string]float64{fromWallet: BtcRoundAmount(maxBalance-amount, 8)})
	}

	request := BtcServerRequest{}
	request.Method = "createrawtransaction"
	request.Id = 1
	request.Params = []interface{}{sendFrom, toSend}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswerCreateRawTransaction{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		//log.Fatalf("Unmarshal: %v", err)
		log.Printf("Unmarshal: %v", err)
		return ""
	}

	return result.Result
}

func BtcSighRaw(Raw string) string {

	if Raw == "" {
		return ""
	}

	type StructData struct {
		Hex string `json:"hex"`
	}

	type ServerAnswerSighRawTransaction struct {
		Error  string     `json:"error"`
		Id     uint       `json:"id"`
		Result StructData `json:"result"`
	}

	request := BtcServerRequest{}
	request.Method = "signrawtransactionwithwallet"
	request.Id = 1
	request.Params = []interface{}{Raw}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswerSighRawTransaction{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return ""
	}

	return result.Result.Hex
}

func BtcSendRaw(RawSigned string) bool {

	if RawSigned == "" {
		return false
	}

	type ServerAnswerSighRawTransaction struct {
		Error  string `json:"error"`
		Id     uint   `json:"id"`
		Result string `json:"result"`
	}

	request := BtcServerRequest{}
	request.Method = "sendrawtransaction"
	request.Id = 1
	request.Params = []interface{}{RawSigned}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswerSighRawTransaction{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		//log.Fatalf("Unmarshal: %v", err)
		log.Printf("%s", body)
		return false
	}

	if len(result.Result) > 0 {
		return true
	} else {
		log.Printf("error")
		return false
	}

}

func BtcSendMoney(from string, to string, amount float64) bool {
	data := BtcListUnspend(0, 99999, []string{from})

	raw := BtcCreateRawTransaction(data, from, to, amount)
	//log.Printf("key : %s", raw)

	if len(raw) > 0 {
		signraw := BtcSighRaw(raw)
		//log.Printf("key sign : %s", signraw)

		if len(signraw) > 0 {
			return BtcSendRaw(signraw)
		}
	}

	return false
}

func BtcSendMoneyMilti(from []string, to []string, amount float64) bool {
	data := BtcListUnspend(0, 99999, from)

	raw := BtcCreateRawTransactionMixerRand(data, from[0], to, amount)
	//log.Printf("key : %s", raw)

	if len(raw) > 0 {
		signraw := BtcSighRaw(raw)
		//log.Printf("key sign : %s", signraw)

		if len(signraw) > 0 {
			return BtcSendRaw(signraw)
		}
	}

	return false
}

func BtcGetNewAddress() string {
	type ServerAnswer struct {
		Error  string `json:"error"`
		Id     uint   `json:"id"`
		Result string `json:"result"`
	}

	request := BtcServerRequest{}
	request.Method = "getnewaddress"
	request.Id = 1
	request.Params = []interface{}{}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	result := ServerAnswer{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return ""
	}

	if len(result.Result) > 0 {
		return result.Result
	} else {
		return ""
	}
}

func BtcGetBalanceAll() float64 {
	type ServerAnswer struct {
		Error  string  `json:"error"`
		Id     uint    `json:"id"`
		Result float64 `json:"result"`
	}

	request := BtcServerRequest{}
	request.Method = "getbalance"
	request.Id = 1
	request.Params = []interface{}{}

	data, _ := json.Marshal(request)

	body := BtcRequestApi(string(data))

	//fmt.Println(string(body))

	result := ServerAnswer{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return 0
	}

	if result.Result > 0 {
		return result.Result
	} else {
		return 0
	}
}
