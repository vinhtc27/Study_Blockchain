package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    0,
	}

	proofOfWork := NewProof(block)
	nonce, hash := proofOfWork.Run()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (block *Block) Serialize() []byte {
	result := bytes.Buffer{}
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	HandleErr(err)

	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	block := &Block{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	HandleErr(err)
	return block
}

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
