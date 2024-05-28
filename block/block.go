package block

import (
	"block_chain/utils"
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	TimeStamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte

	// Nonce game
	Nonce  int64
	Target []byte
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexint(b.TimeStamp), b.PrevHash, b.Data, utils.ToHexint(b.Nonce), b.Target}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func CreateBlock(prevhash, data []byte) *Block {
	block := Block{TimeStamp: time.Now().Unix(), Hash: []byte{}, PrevHash: prevhash, Data: data}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

func GeneisBlock() *Block {
	genesisWords := "Hello Block Chain"
	return CreateBlock([]byte{}, []byte(genesisWords))
}
