package main

import (
	"blockchain/blockchain"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	Blockchain *blockchain.BlockChain
}

func (commandLine *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("add -block BLOCK_DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain")
}

func (commandLine *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		commandLine.printUsage()
		runtime.Goexit()
	}
}

func (commandLine *CommandLine) addBlock(data string) {
	commandLine.Blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (commandLine *CommandLine) printChain() {
	iterator := commandLine.Blockchain.Iterator()

	for {
		block := iterator.Next()
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Block hash: %x\n", block.Hash)

		proofOfWork := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (commandLine *CommandLine) run() {
	commandLine.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	default:
		commandLine.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		commandLine.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		commandLine.printChain()
	}

}

func main() {
	defer os.Exit(0)
	blockChain := blockchain.InitBlockChain()
	defer blockChain.Database.Close()

	commandLine := CommandLine{
		Blockchain: blockChain,
	}

	commandLine.run()

}
