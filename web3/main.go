package main

import (
	_ "crypto/sha256"
)

func main() {
	bc := newBlockChain("14PxkwD8cTpzNAT1PYXRwK4qRNbkBVtgFP")
	cli := CLI{bc}
	cli.Run()

	//通过命令行调用后面的就不用了
	// bc.AddBlock("张三向李四转了50元")
	// bc.AddBlock("李四向王五转了100元")

	// //使用迭代器进行打印
	// it := bc.NewIterator()
	// for {
	// 	block := it.Next()
	// 	fmt.Printf("===========================\n\n")
	// 	fmt.Printf("前区块哈希：%x\n", block.PrevHash)
	// 	fmt.Printf("当前区块哈希：%x\n", block.Hash)
	// 	fmt.Printf("前区块数据：%s\n", block.Data)
	// 	fmt.Printf("Timestamp%d\n", block.Timestamp)
	// 	if len(block.PrevHash) == 0 {
	// 		break
	// 	}
	// }

}
