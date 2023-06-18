package main

import (
	"crypto/sha256"
	"time"
)

// 1定义结构
type Block struct {
	//1版本号
	Version uint64
	//2前区块hash
	PrevHash []byte
	//3merkel根
	MerkelRoot []byte
	//4时间戳
	Timestamp uint64
	//5难度值
	Difficulty uint64
	//6随机数
	Nonce uint64
	//a 当前区块hash正常比特币区块中没有当前区块的哈希，这样写是为了方便，当前区块的哈希不存在在区块链中，存在另外的地方，节省区块的空间
	Hash []byte
	//b 数据
	Data []byte
}

//补充区块字段

//更新计算哈希函数
//优化代码

// uint64-->>[]byte
func Uint64ToByte(uint64) []byte {

	return []byte{}
}

// 2.创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		Timestamp:  uint64(time.Now().Unix()),
		Hash:       []byte{},
		Data:       []byte(data),
	}
	block.SetHash()
	return &block
}

// 3.生成hash
func (block *Block) SetHash() {
	var blockInfo []byte
	//拼装数据
	blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, block.PrevHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Timestamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	blockInfo = append(blockInfo, block.Hash...)
	blockInfo = append(blockInfo, block.Data...)
	//sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
