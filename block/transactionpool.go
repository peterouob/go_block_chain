package block

import (
	"block_chain/constoce"
	"block_chain/transaction"
	"block_chain/utils"
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

type TransactionPool struct {
	PubTx []*transaction.TransAction
}

func (tp *TransactionPool) AddTransaction(tx *transaction.TransAction) {
	tp.PubTx = append(tp.PubTx, tx)
}

func (tp *TransactionPool) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	if err := encoder.Encode(tp); err != nil {
		log.Println("Error to encode transaction pool :", err)
	}
	if err := os.WriteFile(constoce.TransactionPoolFile, content.Bytes(), 0644); err != nil {
		log.Println("Error to write file :", err)
	}
}

func (tp *TransactionPool) Loadfile() error {
	if !utils.FileExsit(constoce.TransactionPoolFile) {
		return nil
	}
	var transactionPool TransactionPool
	fileContent, err := os.ReadFile(constoce.TransactionPoolFile)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	if err := decoder.Decode(&transactionPool); err != nil {
		return err
	}
	tp.PubTx = transactionPool.PubTx
	return nil
}

func CreateTransactionPool() *TransactionPool {
	transactionPool := TransactionPool{}
	if err := transactionPool.Loadfile(); err != nil {
		log.Println("Load file Error :", err)
	}
	return &transactionPool
}

func RemoveTransactionPoolFile() error {
	err := os.Remove(constoce.TransactionPoolFile)
	return err
}
