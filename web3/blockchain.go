package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// 4.引入区块链
type BlockChain struct {
	// block []*Block
	//使用数据库进行替换，从数据库中拿
	db *bolt.DB
	//对最后一个区块的哈希进行存储
	tail []byte
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

func newBlockChain() *BlockChain {
	// genesisBlock := GenesisBlock()
	// return &BlockChain{
	// 	block: []*Block{genesisBlock},
	// }
	var lastHash []byte
	db, err := bolt.Open(blockChainDb, 0600, nil)
	// defer db.Close() error;
	if err != nil {
		log.Panic("打开数据库失败")
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉，我们需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket失败")
			}

			//创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock()

			//3. 写数据
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

// 创建创世区块，在新建区块链时使用
func GenesisBlock() *Block {
	return NewBlock("创世区块建成,很强", []byte{})
}

// 5.添加区块
func (bc *BlockChain) AddBlock(data string) {
	db := bc.db
	lastHash := bc.tail
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Printf("bucket不应该为空,请检查输入文件名或者是否有open失败")
		}

		block := NewBlock(data, lastHash)

		//添加区块到区块链的数据库中
		//hash作为key， block的字节流作为value，尚未实现
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)

		//更新下区块链的tail为当前块的hash
		bc.tail = block.Hash
		return nil
	})
}
