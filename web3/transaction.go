package main

import (
	"crypto/ecdsa"
	"crypto/rand"
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
	privateKey := wallet.Private
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

	bc.SignTransaction(&tx, privateKey)

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

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//todo具体的签名操作

	//创建交易副本
	txCopy := tx.TrimmedCopy()
	//循环遍历inputs，得到input索引的output公钥哈希
	for i, input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		//不要对input进行赋值，这是一个副本，要对txCopy.TXInputs[xx]进行操作，否则无法把pubKeyHash传进来
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PublicKeyHash
		//所需要的三个数据都具备了，开始做哈希处理
		//3. 生成要签名的数据。要签名的数据一定是哈希值
		//a. 我们对每一个input都要签名一次，签名的数据是由当前input引用的output的哈希+当前的outputs（都承载在当前这个txCopy里面）
		//b. 要对这个拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名最终数据。
		txCopy.SetHash()

		//还原，以免影响后面input的签名
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXID
		//4. 执行签名动作得到r,s字节流
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic(err)
		}

		//5. 放到我们所签名的input的Signature中
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXid, input.Index, nil, nil})
	}

	for _, output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}

	return Transaction{tx.TXID, inputs, outputs}
}
