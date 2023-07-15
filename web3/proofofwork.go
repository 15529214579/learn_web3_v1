package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	//a. block
	block *Block
	//b. 目标值
	//一个非常大数
	target *big.Int
}

// 2. 提供创建POW的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//我们指定的难度值，现在是一个string类型，需要进行转换
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr, 16)

	pow.target = &tmpInt
	return &pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//1. 拼装数据（区块的数据，还有不断变化的随机数）
	//2. 做哈希运算
	//3. 与pow中的target进行比较
	//a. 找到了，退出返回
	//b. 没找到，继续找，随机数加1

	var nonce uint64
	//block := pow.block
	var hash [32]byte
	fmt.Print("开始挖矿喽！")

	for {

		// fmt.Printf("hash : %x\r", hash)

		//1. 拼装数据（区块的数据，还有不断变化的随机数）
		blockInfo := pow.PrepareData(nonce)
		//2. 做哈希运算
		hash = sha256.Sum256(blockInfo)
		//3. 与pow中的target进行比较
		tmpInt := big.Int{}
		//将我们得到hash数组转换成一个big.int
		tmpInt.SetBytes(hash[:])

		if tmpInt.Cmp(pow.target) == -1 {
			//a. 找到了，退出返回
			fmt.Printf("挖矿成功！hash : %x, nonce : %d\n", hash, nonce)
			//break
			return hash[:], nonce
		} else {
			//b. 没找到，继续找，随机数加1
			nonce++
		}
	}

}

func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.Timestamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(nonce),
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})
	return blockInfo
}
