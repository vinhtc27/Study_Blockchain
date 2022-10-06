package blockchain

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

type BlockChain struct {
	Blocks []*Block
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

func (blockChain *BlockChain) AddBlock(data string) {
	prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	blockChain.Blocks = append(blockChain.Blocks, newBlock)
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
