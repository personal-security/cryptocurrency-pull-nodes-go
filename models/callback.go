package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type CallbackWallet struct {
	gorm.Model `json:"-"`

	Node    string `json:"node"`
	Network string `json:"network"`
	Address string `json:"address"`
	Url     string `json:"url"`
}

func (callbackWallet *CallbackWallet) Create() (bool, string) {
	if callbackWallet.Address == "" {
		return false, "Error address value"
	}

	if callbackWallet.Url == "" {
		return false, "Error callback url value"
	}

	if callbackWallet.Node == "" {
		return false, "Error node value"
	}

	if callbackWallet.Network == "" {
		return false, "Error network value"
	}

	fundQuery := CallbackWalletGet(callbackWallet.Node, callbackWallet.Network, callbackWallet.Address)
	if fundQuery == nil {
		GetDB().Create(&callbackWallet)
		return true, "Added new callback to db"
	} else {
		fundQuery.Url = callbackWallet.Url
		GetDB().Update(&fundQuery)
		return true, "Update callback in db"
	}

}

func CallbackWalletGet(node string, network string, wallet string) *CallbackWallet {
	callbackWallet := &CallbackWallet{}

	err := GetDB().Where("node = ? AND network = ? AND address = ?", node, network, wallet).Find(&callbackWallet).Limit(1).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return callbackWallet
}
