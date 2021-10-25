package main

import (
	"coin-nodes-rest/app"
	"coin-nodes-rest/controllers"
	"coin-nodes-rest/crons"
	"coin-nodes-rest/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
)

func main() {

	clientHandle := &controllers.ClientHandler{
		Eth: models.EthInitClient(),
		Lru: models.LruInit(),
	}

	//start cron
	c := cron.New()

	c.AddFunc("0 * * * *", crons.CronStatus)
	c.AddFunc("*/5 * * * *", crons.CronCheckCallBack)
	c.AddFunc("*/5 * * * *", func() {
		log.Println("Check coins list (CRON)")
		crons.CronsGetDefaultPrices(clientHandle.Lru)
	})

	c.Start()

	//generate api for builder
	router := mux.NewRouter()

	versionApi1 := "v1"
	prefixQuery1 := "/api/" + versionApi1 + "/"

	//status handlers
	router.HandleFunc(prefixQuery1+"status", controllers.StatusGetNow).Methods("GET")

	router.HandleFunc(prefixQuery1+"list/coins", controllers.CoinsList).Methods("GET")

	router.HandleFunc(prefixQuery1+"prices/coins/list", controllers.CoinsPriceList).Methods("GET")
	router.HandleFunc(prefixQuery1+"prices/coins/list/prepared", controllers.CoinsPriceListPrepared).Methods("GET")
	router.HandleFunc(prefixQuery1+"prices/coins/list/prepared", controllers.CoinsPriceListPreparedFiltered).Methods("POST")
	router.HandleFunc(prefixQuery1+"prices/coins/list/prepared/json", clientHandle.CoinsPriceListPreparedFilteredJson).Methods("POST")
	router.HandleFunc(prefixQuery1+"prices/all", controllers.CoinsPriceAll).Methods("GET")
	router.HandleFunc(prefixQuery1+"prices/current/{coin:[a-z0-9]+}", controllers.CoinsPriceCurrent).Methods("GET")

	router.HandleFunc(prefixQuery1+"callback/node/{type_wallet:[a-z0-9]+}/wallet/{wallet:[a-zA-Z0-9]+}/get", controllers.CallbackGet).Methods("POST")
	router.HandleFunc(prefixQuery1+"callback/node/{type_wallet:[a-z0-9]+}/wallet/{wallet:[a-zA-Z0-9]+}/set", controllers.CallbackSet).Methods("POST")

	// deprecated block
	router.HandleFunc(prefixQuery1+"auth", controllers.Auth).Methods("POST")
	router.HandleFunc(prefixQuery1+"auth/logout", controllers.AuthLogout).Methods("POST")

	router.HandleFunc(prefixQuery1+"nodes/info-all", controllers.NodesInfoAll).Methods("GET")

	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/info", controllers.NodeInfo).Methods("GET")
	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/new/wallet", controllers.WalletGenerate).Methods("GET")
	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/wallet/{wallet:[a-zA-Z0-9]+}", clientHandle.WalletInfo).Methods("GET")
	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/tx/{tx:[a-z0-9]+}", clientHandle.TransactionInfo).Methods("GET")
	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/helper/transaction", clientHandle.TransactionHelper).Methods("POST")

	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/send", controllers.MoneySend).Methods("POST")
	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/multisend", controllers.MoneyMultiSend).Methods("POST")
	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/broadcast/raw", clientHandle.BroadcastRaw).Methods("POST")

	router.HandleFunc(prefixQuery1+"node/{type_wallet:[a-z0-9]+}/tools/restart", controllers.ToolsNodeRestart).Methods("GET")

	router.HandleFunc(prefixQuery1+"system/backup", controllers.ToolsBackup).Methods("GET")
	router.HandleFunc(prefixQuery1+"system/restore", controllers.ToolsRestore).Methods("POST")

	router.HandleFunc(prefixQuery1+"system/restart-all", controllers.ToolsAllNodesRestart).Methods("GET")

	router.Use(app.KeyAuthentication) //attach API auth middleware

	router.NotFoundHandler = http.HandlerFunc(app.NotFoundHandler)

	type ConfigNetwork struct {
		Ip   string `json:"ip"`
		Port string `json:"port"`
	}

	var configNetwork ConfigNetwork

	file, _ := ioutil.ReadFile("config/network.json")

	json.Unmarshal([]byte(file), &configNetwork)

	if configNetwork.Ip == "" || configNetwork.Port == "" {
		fmt.Println("Error network config")
		return
	}

	log.Printf("\nServer start on %s:%s \n", configNetwork.Ip, configNetwork.Port)

	err := http.ListenAndServe(configNetwork.Ip+":"+configNetwork.Port, router)
	if err != nil {
		fmt.Println(err)
	}
}
