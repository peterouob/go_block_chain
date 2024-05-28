package block

import (
	"block_chain/constoce"
	"block_chain/utils"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	// 左移越少難度越大,和difficult成正比
	target.Lsh(target, uint(256-constoce.Difficult))
	return target.Bytes()
}

// 每次都會產生新的nano
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexint(b.TimeStamp),
		b.PrevHash,
		b.Data,
		utils.ToHexint(nonce),
		b.Target,
	}, []byte{})
	return data
}

func (b *Block) FindNonce() int64 {
	var (
		intHash   big.Int
		intTarget big.Int
		hash      [32]byte
		nonce     int64
	)
	nonce = 0
	intTarget.SetBytes(b.Target)
	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			// 直到nonce小於目標難度
			break
		} else {
			nonce++
		}
	}
	return nonce
}

func (b *Block) ValidPoW() bool {
	var (
		intHash   big.Int
		intTarget big.Int
		hash      [32]byte
	)
	intTarget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(&intTarget) == -1
}
