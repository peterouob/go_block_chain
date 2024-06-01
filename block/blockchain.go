package block

import (
	"block_chain/transaction"
	"encoding/hex"
	"log"
)

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(txs []*transaction.TransAction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func CreateBlockChain() *BlockChain {
	blockChain := BlockChain{}
	blockChain.Blocks = append(blockChain.Blocks, GeneisBlock())
	return &blockChain
}

func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.TransAction {
	// 包含指定地址的可用交易信息切片
	var unSpentTxs []transaction.TransAction
	// Key 值為交易訊息ID value為Output在該交易中的序號，此為紀錄邊例區塊鏈時已經被使用的交易訊息Output
	spentTxs := make(map[string][]int)
	// 遍歷交易區塊訊息
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]
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

func (bc *BlockChain) Mine(txs []*transaction.TransAction) {
	bc.AddBlock(txs)
}
