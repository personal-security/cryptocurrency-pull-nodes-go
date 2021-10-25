package controllers

import (
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/flyaways/golang-lru"
	"github.com/gorilla/schema"
)

type ClientHandler struct {
	Eth *ethclient.Client
	Lru *lru.Cache
}

var decoder = schema.NewDecoder()
