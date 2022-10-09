package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	lastHash := []byte{}

	options := badger.DefaultOptions(dbPath)
	database, err := badger.Open(options)
	HandleErr(err)

	err = database.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("lastHash"))
		if err == badger.ErrKeyNotFound {
			genesis := Genesis()
			fmt.Println("Genesis proved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleErr(err)

			err = txn.Set([]byte("lastHash"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lastHash"))
			HandleErr(err)

			lastHash, err = item.ValueCopy(lastHash)
			return err
		}

	})
	HandleErr(err)

	blockchain := &BlockChain{
		LastHash: lastHash,
		Database: database,
	}

	return blockchain
}

func (blockChain *BlockChain) AddBlock(data string) {
	lastHash := []byte{}

	err := blockChain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		HandleErr(err)

		lastHash, err = item.ValueCopy(lastHash)
		return err
	})
	HandleErr(err)

	newBlock := CreateBlock(data, lastHash)

	err = blockChain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleErr(err)

		err = txn.Set([]byte("lastHash"), newBlock.Hash)

		blockChain.LastHash = newBlock.Hash

		return err
	})
	HandleErr(err)
}

func (blockChain *BlockChain) Iterator() *BlockChainIterator {
	iterator := &BlockChainIterator{
		CurrentHash: blockChain.LastHash,
		Database:    blockChain.Database,
	}

	return iterator
}

func (iterator *BlockChainIterator) Next() *Block {
	block := &Block{}

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		HandleErr(err)

		encodedBlock := []byte{}

		encodedBlock, err = item.ValueCopy(encodedBlock)

		block = Deserialize(encodedBlock)

		return err
	})
	HandleErr(err)

	iterator.CurrentHash = block.PrevHash

	return block
}
