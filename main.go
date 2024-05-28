package main

import (
	"block_chain/block"
	"fmt"
	"time"
)

func main() {
	blockchain := block.CreateBlockChain()
	time.Sleep(time.Second)
	blockchain.AddBlock("Hello i am peter Lin ")
	time.Sleep(time.Second)
	blockchain.AddBlock("Hello i am Defer  ")
	time.Sleep(time.Second)
	blockchain.AddBlock("Hello i am peterdefer ")
	time.Sleep(time.Second)

	for _, block := range blockchain.Blocks {
		fmt.Printf("Time stamp %d\n", block.TimeStamp)
		fmt.Printf("Hash %x\n", block.Hash)
		fmt.Printf("Prev Hash %x\n", block.PrevHash)
		fmt.Printf("Data %s\n", block.Data)
	}
}
