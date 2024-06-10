package block

import (
	"fmt"
	"log"
)

func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx)
	if candidateBlock.ValidPoW() {
		bc.AddBlock(candidateBlock)
		if err := RemoveTransactionPoolFile(); err != nil {
			log.Println("Remove transaction pool file :", err)
			return
		}
	} else {
		fmt.Println("Block has invalid nonce. ")
		return
	}
}
