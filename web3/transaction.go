package main

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
