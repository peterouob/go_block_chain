package transaction

import (
	"block_chain/constoce"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type TransAction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// 你可能对于为什么要记录一组TxOutput有疑惑，这是因为在寻找到了足够多的未使用的TxOutput（后面全部简称UTXO）后，其资产总量可能大于我们本次交易的转账总量，我们可以将找零计入本次的TxOutput中，设置其流入方向就是本次交易的Sender（一定要好好理解！），这样就实现了找零。

func (tx *TransAction) TxHash() []byte {
	var (
		encoded bytes.Buffer
		hash    [32]byte
	)
	encoder := gob.NewEncoder(&encoded)
	if err := encoder.Encode(tx); err != nil {
		log.Println("Cannnot encoding the struct from trans action :", err)
		return nil
	}
	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *TransAction) SetID() {
	tx.ID = tx.TxHash()
}

func BaseTx(toaddress []byte) *TransAction {
	txIn := TxInput{[]byte{}, -1, []byte{}}
	txOut := TxOutput{Value: constoce.InitCoin, ToAddress: toaddress}
	tx := TransAction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

func (tx *TransAction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutId == -1
}
