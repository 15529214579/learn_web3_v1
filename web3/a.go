package main

import (
	_ "crypto/sha256"
	"fmt"
)

func main() {
	bc := newBlockChain()
	bc.AddBlock("班长向班花转了50元")
	bc.AddBlock("班长向班花转了100元")
	for i, block := range bc.block {
		fmt.Printf("区块高度：%d\n", i)
		fmt.Printf("前区块哈希：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("前区块数据：%s\n", block.Data)
		fmt.Printf("Timestamp%d\n", block.Timestamp)
		// fmt.Println("maxuefei test block")
	}

}
