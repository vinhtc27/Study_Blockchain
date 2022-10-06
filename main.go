package main

import (
	"fmt"
	"strconv"
	"time"
	"vinhtc-chain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	before := time.Now()
	for index := 1; index <= 100; index++ {
		chain.AddBlock("Block #" + strconv.Itoa(index))
		fmt.Printf("Previous hash: %x\n", chain.Blocks[index].PrevHash)
		fmt.Printf("Data in block: %s\n", chain.Blocks[index].Data)
		fmt.Printf("Block hash: %x\n", chain.Blocks[index].Hash)

		proofOfWork := blockchain.NewProof(chain.Blocks[index])
		fmt.Printf("Proof of work: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println()
	}
	after := time.Now()
	fmt.Println(after.Sub(before).Seconds())
}
