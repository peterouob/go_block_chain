package main

import (
	"block_chain/block"
	"block_chain/transaction"
	"fmt"
)

func main() {
	txPool := make([]*transaction.TransAction, 0)
	var tempTx *transaction.TransAction
	var ok bool
	var property int

	chain := block.CreateBlockChain()

	property, _ = chain.FindUTXOs([]byte("Peter Lin"))
	fmt.Println("Balance of Peter Lin:", property)

	tempTx, ok = chain.CreateTransaction([]byte("Peter Lin"), []byte("Aaron"), 100)

	if ok {
		txPool = append(txPool, tempTx)
	}

	chain.Mine(txPool)

	txPool = make([]*transaction.TransAction, 0)
	property, _ = chain.FindUTXOs([]byte("Aaron"))
	fmt.Println("Balance of Aaron:", property)
	tempTx, ok = chain.CreateTransaction([]byte("Aaron"), []byte("Alisa"), 200)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("Aaron"), []byte("Alisa"), 100)
	if ok {
		txPool = append(txPool, tempTx)
	}

	tempTx, ok = chain.CreateTransaction([]byte("Peter Lin"), []byte("Alisa"), 50)
	if ok {
		txPool = append(txPool, tempTx)
	}
	chain.Mine(txPool)
	txPool = make([]*transaction.TransAction, 0)
	property, _ = chain.FindUTXOs([]byte("Peter Lin"))
	fmt.Println("Balance of Peter Lin:", property)

	property, _ = chain.FindUTXOs([]byte("Peter Aaron"))
	fmt.Println("Balance of Aaron:", property)

	property, _ = chain.FindUTXOs([]byte("Alisa"))
	fmt.Println("Balance of Alisa:", property)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.TimeStamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Println("Proof of Work validation: ", block.ValidPoW())
	}
}
