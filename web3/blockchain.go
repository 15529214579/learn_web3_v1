package main

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
	return NewBlock("创世区块建成,很强", []byte{})
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
