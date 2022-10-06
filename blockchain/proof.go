package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

var _difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-_difficulty))

	return &ProofOfWork{
		Block:  block,
		Target: target,
	}
}

func (proofOfWork *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			proofOfWork.Block.PrevHash,
			proofOfWork.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(_difficulty)),
		},
		[]byte{},
	)
}

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	intHash := big.Int{}
	hash := [32]byte{}

	nonce := 0

	for nonce < math.MaxInt64 {
		data := proofOfWork.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(proofOfWork.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (proofOfWork *ProofOfWork) Validate() bool {
	intHash := big.Int{}
	data := proofOfWork.InitData(proofOfWork.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(proofOfWork.Target) == -1

}

func ToHex(number int64) []byte {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, number)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
