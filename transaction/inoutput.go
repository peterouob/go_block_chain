package transaction

import (
	"bytes"
)

// 轉出
type TxOutput struct {
	Value     int
	ToAddress []byte
}

// 收入
type TxInput struct {
	TxID        []byte
	OutId       int
	FromAddress []byte
}

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
