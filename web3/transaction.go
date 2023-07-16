package main

import (
	"crypto/sha256"
	"log"

	"github.com/bytedance/sonic"
)

const reward = 12.5

type Transaction struct {
	TXID      []byte     //交易id
	TXInputs  []TXInput  //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}

// 交易输入
type TXInput struct {
	TXid  []byte //引用的输出交易id
	Index int64  //引用的输出索引值
	Sig   string //解锁脚本
}

// 交易输出
type TXOutput struct {
	value         float64 //转账的金额
	PublicKeyHash string  //锁定脚本，先用公钥哈希来模拟（后期要替换成脚本黑盒）
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
	//没有index，没有input交易id,解锁脚本先用data代替
	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}

	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()

	return &tx
}
