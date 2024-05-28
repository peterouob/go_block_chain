package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	TimeStamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{ToHexint(b.TimeStamp), b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func ToHexint(num int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		log.Println("Error to translate num to []byte :", err)
		return nil
	}
	return buff.Bytes()
}

func CreateBlock(prevhash, data []byte) *Block {
	block := Block{TimeStamp: time.Now().Unix(), Hash: []byte{}, PrevHash: prevhash, Data: data}
	block.SetHash()
	return &block
}

func GeneisBlock() *Block {
	genesisWords := "Hello Block Chain"
	return CreateBlock([]byte{}, []byte(genesisWords))
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

func main() {
	blockchain := CreateBlockChain()
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
