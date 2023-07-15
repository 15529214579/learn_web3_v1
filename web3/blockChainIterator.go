package main

import (
	"log"

	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	db             *bolt.DB
	curHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.db,
		bc.tail,
	}
}

func (it *BlockChainIterator) Next() *Block {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器打开bucket时为空，error，请检查")
		}

		blockTmp := bucket.Get(it.curHashPointer)
		block = DeSerialize(blockTmp)
		it.curHashPointer = block.PrevHash
		return nil
	})
	return &block
}
