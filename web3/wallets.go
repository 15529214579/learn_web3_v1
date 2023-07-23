package main

//定一个wallets结构，它保存所有的wallet以及它的地址

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	wallet := NewWallet()
	address := wallet.NewAddress()

	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.WalletsMap[address] = wallet

	return &wallets
}
