package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"

	"github.com/bytedance/sonic"
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
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// 2.创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: 0, //先随便设置一个值
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	block.Serialize()
	return &block
}

// 还没有做，这里需要将整个区块链进行hash操作
func (block *Block) toByte() []byte {
	//TODO maxuefei
	return []byte{}
}

func (block *Block) Serialize() []byte {
	bytes, err := sonic.Marshal(block)
	if err != nil {
		log.Print("序列化出错")
	}
	return bytes
}

func DeSerialize(data []byte) Block {
	var block2 Block
	sonic.Unmarshal(data, &block2)
	return block2
}

// // 3.生成hash
// func (block *Block) SetHash() {
// 	var blockInfo []byte
// 	//拼装数据
// 	/*blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
// 	blockInfo = append(blockInfo, block.PrevHash...)
// 	blockInfo = append(blockInfo, block.MerkelRoot...)
// 	blockInfo = append(blockInfo, Uint64ToByte(block.Timestamp)...)
// 	blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
// 	blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
// 	blockInfo = append(blockInfo, block.Hash...)
// 	blockInfo = append(blockInfo, block.Data...)*/
// 	//sha256
// 	temp := [][]byte{
// 		Uint64ToByte(block.Version),
// 		block.PrevHash,
// 		block.MerkelRoot,
// 		Uint64ToByte(block.Timestamp),
// 		Uint64ToByte(block.Difficulty),
// 		Uint64ToByte(block.Nonce),
// 		block.Hash,
// 		block.Data,
// 	}
// 	//2维数组转化成一维数组
// 	blockInfo = bytes.Join(temp, []byte{})
// 	hash := sha256.Sum256(blockInfo)
// 	block.Hash = hash[:]
// }
