package main

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

func (cli *CLI) AddBlock(data string) {
	// cli.bc.AddBlock(data)
	//todo maxuefei，这里要的是交易，暂时还没发提供
	fmt.Printf("添加区块成功！\n")
}

func (cli *CLI) PrinBlockChain() {
	bc := cli.bc
	//创建迭代器
	it := bc.NewIterator()

	//调用迭代器，返回我们的每一个区块数据
	for {
		//返回区块，左移
		block := it.Next()

		fmt.Printf("===========================\n\n")
		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
		fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %d\n", block.Timestamp)
		fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
		fmt.Printf("随机数 : %d\n", block.Nonce)
		fmt.Printf("当前区块哈希值: %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].PubKey)

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束！")
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	// 1.对地址做校验，如果是非法地址就不能获取余额
	// 2.通过地址获取公钥哈希
	publicKeyHash := GetPubKeyFromAddress(address)
	utxos := cli.bc.findUTXOs(publicKeyHash)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("\"%s\"的余额为：%f\n", address, total)

}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	fmt.Printf("from : %s\n", from)
	fmt.Printf("to : %s\n", to)
	fmt.Printf("amount : %f\n", amount)
	fmt.Printf("miner : %s\n", miner)
	fmt.Printf("data : %s\n", data)

	//创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	txs := []*Transaction{coinbase}
	//创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if txs != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效的交易!\n")
	}
	//添加到区块
	cli.bc.AddBlock(txs)
	fmt.Printf("转账成功!")
}

func (cli *CLI) NewWallet() {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址：%s\n", address)
}

func (cli *CLI) ListAddresses() {
	ws := NewWallets() //只是load本地的数据
	addresses := ws.ListAddresses()
	for _, address := range addresses {
		fmt.Printf("地址:%s\n", address)
	}
}

func GetPubKeyFromAddress(address string) []byte {
	addrDec := base58.Decode(address)
	lenAddr := len(addrDec)
	pubKeyHash := addrDec[1 : lenAddr-4]

	return pubKeyHash
}
