package models

// in development
type TransactionCoin struct {
}

type Coin interface {
	GenerateNewAddress() string
	GetWalletBalance(wallet string) float64
	SendMoney(walletFrom string, walletTo string, amount float64, fee float64) bool

	GatWalletTransactiots(wallet string) []TransactionCoin
	GetTransactionInfo(tx string) TransactionCoin

	GetPoolBalance() float64

	GetFee() float64

	RestartService() bool
}

type AbstractCoin struct{ Coin }

// work code here (for tests)
type TransactionInfo struct {
	Tx string `json:"tx"`
}

type WalletInfo struct {
	Address           string            `json:"address"`
	Amount            float64           `json:"amount"`
	AmountUnconfirmed float64           `json:"amountUnconfirmed"`
	RecomendedFee     float64           `json:"recomendedFee"`
	Transactions      []TransactionInfo `json:"transactions"`
}
