package crons

import (
	"log"
	"net/http"
	"time"

	lru "github.com/flyaways/golang-lru"
	coingecko "github.com/superoo7/go-gecko/v3"
)

func CronsGetDefaultPrices(cache *lru.Cache) {

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	list, err := cg.CoinsList()
	if err != nil {
		log.Println(err)
	}

	cache.Add("coins_list", list)
}
