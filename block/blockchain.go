package block

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
}

func CreateBlockChain() *BlockChain {
	blockChain := BlockChain{}
	blockChain.Blocks = append(blockChain.Blocks, GeneisBlock())
	return &blockChain
}
