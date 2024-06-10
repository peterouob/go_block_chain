package block

import (
	"block_chain/constoce"
	"block_chain/transaction"
	"block_chain/utils"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"

	"github.com/dgraph-io/badger"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}
	return &iterator
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block
	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		if err != nil {
			log.Println("Get iterator current hash error :", err)
			return err
		}
		err = item.Value(func(val []byte) error {
			block = block.DeSerialize(val)
			return nil
		})
		if err != nil {
			log.Println("Block chain vdeserialize value error :", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("View database error :", err)
		return nil
	}
	iterator.CurrentHash = block.PrevHash
	return block
}

func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		if err != nil {
			log.Println("get ogprevhash Error")
			return err
		}
		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})
		if err != nil {
			log.Println("Error to get item value :", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Database View error :", err)
		return nil
	}
	return ogprevhash
}
func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte
	if utils.FileExsit(constoce.BCFile) {
		fmt.Println("block chain is already exist")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constoce.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Println("Have Error on Open DB :", err)
		return nil
	}

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GeneisBlock(address)
		fmt.Println("Genesis Create")
		if err := txn.Set(genesis.Hash, genesis.Serialize()); err != nil {
			log.Println("Error on Set genesis :", err)
			return err
		}
		if err := txn.Set([]byte("lh"), genesis.Hash); err != nil {
			log.Println("Error on Set genesis :", err)
			return err
		}
		if err := txn.Set([]byte("ogprevhash"), genesis.PrevHash); err != nil {
			log.Println("Error on Set genesis :", err)
			return err
		}
		lastHash = genesis.Hash
		return nil
	})
	if err != nil {
		log.Println("Error on Update :", err)
	}
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain() *BlockChain {
	if utils.FileExsit(constoce.BCPath) == false {
		log.Println("No blockchain found, please create one first")
		return nil
	}
	var lastHash []byte
	opts := badger.DefaultOptions(constoce.BCPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Println("Cannot open the DB :", err)
		return nil
	}
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Println("Get Item from ln error :", err)
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		if err != nil {
			log.Println("Error to get item Value :", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Error to update :", err)
		return nil
	}
	chain := BlockChain{lastHash, db}
	return &chain
}

func (bc *BlockChain) AddBlock(newBlock *Block) {
	var lastHash []byte
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Println("Error to Get :", err)
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		if err != nil {
			log.Println("Error to Get Val :", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Error to View :", err)
	}
	if !bytes.Equal(newBlock.PrevHash, lastHash) {
		log.Println("THIS BLOCK IS OUT OF AGE")
		runtime.Goexit()
	}

	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Println("Set new Block error :", err)
			return err
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	if err != nil {
		log.Println("Update Error :", err)
	}
}

func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.TransAction {
	// 包含指定地址的可用交易信息切片
	var unSpentTxs []transaction.TransAction
	// Key 值為交易訊息ID value為Output在該交易中的序號，此為紀錄邊例區塊鏈時已經被使用的交易訊息Output
	spentTxs := make(map[string][]int)

	iter := bc.Iterator()
all:
	// 遍歷交易區塊訊息
	for {
		block := iter.Next()
		for _, tx := range block.Transaction {
			txID := hex.EncodeToString(tx.ID)
		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spenOut := range spentTxs[txID] {
						// 代表交易訊息已經被使用過，跳過
						if spenOut == outIdx {
							// 繼續循環迭代
							continue IterOutputs
						}
					}
				}
				// 確認是否和to address一樣，正確就是我們要找的交易訊息
				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}
			// 檢查是否沒有input，如果不是就檢查當前交易訊息的input是否含目標地址，有的話就將Output訊息加入到spentTxs
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutId)
					}
				}
			}
		}
		if bytes.Equal(block.PrevHash, bc.BackOgPrevHash()) {
			break all
		}
	}
	return unSpentTxs
}

func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accountMount := 0
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIDx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				accountMount += out.Value
				unspentOuts[txID] = outIDx
				continue Work
			}
		}
	}
	return accountMount, unspentOuts
}

func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += amount
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.TransAction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Println("Not enough coins! coin:", acc)
		return &transaction.TransAction{}, false
	}
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Println(err)
			return &transaction.TransAction{}, false
		}
		input := transaction.TxInput{TxID: txID, OutId: outidx, FromAddress: from}
		inputs = append(inputs, input)
	}
	outputs = append(outputs, transaction.TxOutput{Value: amount, ToAddress: to})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, ToAddress: from})
	}
	tx := transaction.TransAction{ID: nil, Inputs: inputs, Outputs: outputs}
	tx.SetID()
	return &tx, true
}
