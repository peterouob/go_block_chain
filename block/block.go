package block

import (
	"block_chain/transaction"
	"block_chain/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
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

	Transaction []*transaction.TransAction
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexint(b.TimeStamp), b.PrevHash, b.Data, utils.ToHexint(b.Nonce), b.Target, b.BackTrasactionSummary()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func CreateBlock(prevhash []byte, txs []*transaction.TransAction) *Block {
	block := Block{TimeStamp: time.Now().Unix(), Hash: []byte{}, PrevHash: prevhash, Data: []byte{}, Nonce: 0, Target: []byte{}, Transaction: txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

func GeneisBlock(address []byte) *Block {
	tx := transaction.BaseTx(address)
	log.Printf("Peter Lin Get Init Coin name :%s ,coin: %d \n", string(tx.Outputs[0].ToAddress), tx.Outputs[0].Value)
	genesis := CreateBlock([]byte{}, []*transaction.TransAction{tx})
	genesis.SetHash()
	return genesis
}

// 由於缺少第三分認證，因次是否有足夠餘額會從上次交易看
func (b *Block) BackTrasactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transaction {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	if err := encoder.Encode(b); err != nil {
		return nil
	}
	return res.Bytes()
}

func (b *Block) DeSerialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		return nil
	}
	return &block
}
