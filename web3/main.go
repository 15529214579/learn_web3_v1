package main

import (
	_ "crypto/sha256"
)

func main() {
	bc := newBlockChain()
	bc.AddBlock("班长向班花转了50元")
	bc.AddBlock("班长向班花转了100元")
	/*
		//改成数据库之后打印不能这样打了，需要从数据库中进行读取
		for i, block := range bc.block {
			fmt.Printf("区块高度：%d\n", i)
			fmt.Printf("前区块哈希：%x\n", block.PrevHash)
			fmt.Printf("当前区块哈希：%x\n", block.Hash)
			fmt.Printf("前区块数据：%s\n", block.Data)
			fmt.Printf("Timestamp%d\n", block.Timestamp)
			// fmt.Println("maxuefei test block")
		}*/

}
