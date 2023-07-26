package main

import (
	"bytes"
	"crypto/ecdsa"
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

func newBlockChain(address string) *BlockChain {
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
			genesisBlock := GenesisBlock(address)

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
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "创世区块建成,很强")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 5.添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//校验是否为合法交易
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Printf("矿工发现无效交易!")
			return
		}
	}

	db := bc.db
	lastHash := bc.tail
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Printf("bucket不应该为空,请检查输入文件名或者是否有open失败")
		}

		block := NewBlock(txs, lastHash)

		//添加区块到区块链的数据库中
		//hash作为key， block的字节流作为value，尚未实现
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)

		//更新下区块链的tail为当前块的hash
		bc.tail = block.Hash
		return nil
	})
}

// 这个是查询余额的时候调用的,找到当前地址的所有utxo
func (bc *BlockChain) findUTXOs(publicKeyHash []byte) []TXOutput {
	var UTXO []TXOutput
	txs := bc.FindUTXOTransactions(publicKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(output.PublicKeyHash, publicKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

func (bc *BlockChain) FindNeedUTXOs(publicKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxos集合
	utxos := make(map[string][]uint64)
	var calc float64

	txs := bc.FindUTXOTransactions(publicKeyHash)

	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			if bytes.Equal(publicKeyHash, output.PublicKeyHash) {

				utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
				calc += output.Value

				if calc >= amount {
					//break
					fmt.Printf("找到了满足的金额：%f\n", calc)
					return utxos, calc
				}
			} else {
				fmt.Printf("不满足转账金额,当前总额：%f， 目标金额: %f\n", calc, amount)
			}
		}
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(publicKeyHash []byte) []*Transaction {
	var txs []*Transaction
	//我们定义一个map来保存消费过的output，key是这个output的交易id，value是这个交易中索引的数组
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	//创建迭代器
	it := bc.NewIterator()

	for {
		//1.遍历区块
		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {

		OUTPUT:
			//3. 遍历output，找到和自己相关的utxo(在添加output之前检查一下是否已经消耗过)
			//	i : 0, 1, 2, 3
			for i, output := range tx.TXOutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						//[]int64{0, 1} , j : 0, 1
						if int64(i) == j {
							//当前准备添加output已经消耗过了，不要再加了
							continue OUTPUT
						}
					}
				}

				//这个output和我们目标的地址相同，满足条件，加到返回UTXO数组中
				if bytes.Equal(output.PublicKeyHash, publicKeyHash) {
					txs = append(txs, tx)
				} else {
				}
			}

			//如果当前交易是挖矿交易的话，那么不做遍历，直接跳过

			if !tx.IsCoinbase() {
				//4. 遍历input，找到自己花费过的utxo的集合(把自己消耗过的标示出来)
				for _, input := range tx.TXInputs {
					//判断一下当前这个input和目标（李四）是否一致，如果相同，说明这个是李四消耗过的output,就加进来
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, publicKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			} else {
				//fmt.Printf("这是coinbase，不做input遍历！")
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return txs
}

func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	it := bc.NewIterator()

	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历结束!\n")
			break
		}
	}
	return Transaction{}, nil
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1. 根据inputs来找，有多少input, 就遍历多少次
	//2. 找到目标交易，（根据TXid来找）
	//3. 添加到prevTXs里面
	for _, input := range tx.TXInputs {
		//根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.TXid)

		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx
	}

	tx.Sign(privateKey, prevTXs)
}
