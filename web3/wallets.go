package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bytedance/sonic"
)

const walletFile = "wallet.dat"

//定一个wallets结构，它保存所有的wallet以及它的地址

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	//只是创建钱包s，具体的添加地址对应哪个钱包放到CreateWallet中
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	ws.WalletsMap[address] = wallet
	ws.saveToFile()
	//返回是地址,然后将钱包s写到文件中
	return address
}

func (ws *Wallets) saveToFile() {
	bytes, err := sonic.Marshal(ws)
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile(walletFile, bytes, 0600)
}
func (ws *Wallets) loadFile() {
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		fmt.Printf("本地钱包文件不存在")
		return
	}
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wsLocalFile Wallets
	sonic.Unmarshal(content, &wsLocalFile)
	ws.WalletsMap = wsLocalFile.WalletsMap
}

func (ws *Wallets) ListAddresses() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}
