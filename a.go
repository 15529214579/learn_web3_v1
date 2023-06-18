package main

import (
	"crypto/sha256"
	"fmt"
)

// 1定义结构
type Block struct {
	//前区块hash
	PrevHash []byte
	//当前区块hash
	Hash []byte
	//数据
	Data []byte
}

// 2.创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

// 3.生成hash
func (block *Block) SetHash() {
	//拼装数据
	blockInfo := append(block.PrevHash, block.Data...)
	//sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

// 4.引入区块链
type BlockChain struct {
	block []*Block
}

func newBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		block: []*Block{genesisBlock},
	}
}

// 创建创世区块，在新建区块链时使用
func GenesisBlock() *Block {
	return NewBlock("创世区块建成,捞牛逼了", []byte{})
}

// 5.添加区块
func (bc *BlockChain) AddBlock(data string) {
	//获取最后一个区块
	lastBlock := bc.block[len(bc.block)-1]
	prevHash := lastBlock.Hash
	//创建区块
	block := NewBlock(data, prevHash)
	bc.block = append(bc.block, block)
	//添加到区块链数组中
}

// 6.重构代码
func main() {
	bc := newBlockChain()
	bc.AddBlock("班长向班花转了50元")
	bc.AddBlock("班长向班花转了100元")
	for i, block := range bc.block {
		fmt.Printf("区块高度：%d\n", i)
		fmt.Printf("前区块哈希：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("前区块数据：%s\n", block.Data)
		// fmt.Println("maxuefei test block")
	}

}
