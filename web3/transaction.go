package main

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
)

const reward = 60

type Transaction struct {
	TXID      []byte     //交易id
	TXInputs  []TXInput  //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}

// 交易输入
type TXInput struct {
	TXid  []byte //引用的输出交易id
	Index int64  //引用的输出索引值
	// Sig   string //解锁脚本
	Signature []byte //rs拼接成的数字签名
	PubKey    []byte //存储的xy拼接后的数据

}

// 交易输出
type TXOutput struct {
	Value         float64 //转账的金额
	PublicKeyHash []byte  //收款方的公钥哈希
}

// 设置交易id
func (tx *Transaction) SetHash() {
	bytes, err := sonic.Marshal(tx)
	if err != nil {
		log.Panic("SetHash error 序列化出错，请检查交易和序列化函数")
	}

	hash := sha256.Sum256(bytes)
	tx.TXID = hash[:]
}

// 创建创世区块
func NewCoinbaseTX(address string, data string) *Transaction {
	//矿工挖矿时不需要指定签名
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	output := NewTXOutput(reward, address)

	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	tx.SetHash()

	return &tx
}

// 实现一个函数，判断当前的交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	//1. 交易input只有一个
	//2. 交易id为空
	//3. 交易的index 为 -1
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}

	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//先打开钱包
	ws := NewWallets()

	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("打开钱包地址失效,交易创建失败")
		return nil
	}

	pubkey := wallet.PubKey
	//privateKey := wallet.Private
	pubkeyHash := HashPubKey(pubkey)

	//todo maxuefei FindNeedUTXOs也需要同步修改
	utxos, resValue := bc.FindNeedUTXOs(pubkeyHash, amount)
	if resValue < amount {
		fmt.Printf("余额不足, 交易失败！")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//创建交易输入，将utxo转化成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), nil, pubkey}
			inputs = append(inputs, input)
		}
	}
	//创建交易输出
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)

	//找零
	if resValue > amount {
		// outputs = append(outputs, TXOutput{resValue - amount, from})
		output = NewTXOutput(resValue-amount, from)
		outputs = append(outputs, *output)
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}

// 由于现在存储的字段是地址的公钥哈希，所以无法直接创建TXOutput，
// 为了能够得到公钥哈希，我们需要处理一下，写一个Lock函数
func (output *TXOutput) Lock(address string) {
	output.PublicKeyHash = GetPubKeyFromAddress(address)
}

// 给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}

	output.Lock(address)
	return &output
}
