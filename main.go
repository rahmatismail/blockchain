package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

type (
	Block struct {
		Timestamp    int64
		Data         []byte
		Hash         []byte
		PreviousHash []byte
		nonce        int
	}

	BlockChain struct {
		blocks []*Block
	}

	ProofOfWork struct {
		block  *Block
		target *big.Int
	}
)

const (
	maxNonce = math.MaxInt64
)

var (
	difficulty = 0 // increase every 1000 blocks
)

// NewProofOfWork for creating proof of work object
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	return &ProofOfWork{
		block:  b,
		target: target,
	}
}

// merge block header
// prepare
// prev hash
// data
// timestamp
// difficulty
// nonce
func (pow *ProofOfWork) prepareHash(nonce int) (mergedData []byte) {
	mergedData = bytes.Join([][]byte{
		pow.block.PreviousHash,
		pow.block.Data,
		[]byte(strconv.FormatInt(pow.block.Timestamp, 16)),
		[]byte(strconv.FormatInt(int64(difficulty), 16)),
		[]byte(strconv.FormatInt(int64(nonce), 16)),
	}, []byte{})

	return
}

// Run will will return nonce and desired block hash based on difficulty
// prepare hash based on block header data and nonce
// calculate hash use sha256
// compare generated hash. if less than target then hash accepted
// repeat prepare hash if requirement above doesn't meet. Increment nonce
func (pow *ProofOfWork) Run() (nonce int, hash []byte) {
	var hashInt big.Int

	fmt.Printf("mining block containing -> %s <- difficulty: %d\n", pow.block.Data, difficulty)
	for nonce < maxNonce {
		mergedData := pow.prepareHash(nonce)
		hash32 := sha256.Sum256(mergedData)
		hash = hash32[:]

		fmt.Printf("\r%x - difficulty:%d", hash, nonce)
		hashInt.SetBytes(hash)

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")

	return
}

// Validate proof of work hash
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	mergedData := pow.prepareHash(pow.block.nonce)
	hash32 := sha256.Sum256(mergedData)

	hashInt.SetBytes(hash32[:])

	return hashInt.Cmp(pow.target) == -1
}

// NewBlock create block object
func NewBlock(data []byte, previousHash []byte) *Block {

	block := &Block{
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: previousHash,
	}

	pow := NewProofOfWork(block)

	block.nonce, block.Hash = pow.Run()
	return block
}

// NewFirstBlock construct first block
func NewFirstBlock() *Block {
	return NewBlock([]byte("First block"), []byte{})
}

// NewBlockChain create blockchain object
func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{
			NewFirstBlock(),
		},
	}
}

// addBlock append new block on block chain
func (bc *BlockChain) addBlock(data string) {
	prevHash := bc.blocks[len(bc.blocks)-1].Hash

	bc.blocks = append(bc.blocks, NewBlock([]byte(data), prevHash))
}

func main() {
	bc := NewBlockChain()

	for i := 0; i < 3; i++ {
		if i%100 == 0 {
			difficulty++
		}
		bc.addBlock(fmt.Sprintf("this is the block[%d]", i))
	}

	for _, b := range bc.blocks {
		fmt.Println("=================")
		fmt.Printf("hash    : %x\n", b.Hash)
		fmt.Printf("time    : %x\n", b.Timestamp)
		fmt.Printf("Prehash : %x\n", b.PreviousHash)
		fmt.Printf("Data    : %s\n", b.Data)

		pow := NewProofOfWork(b)
		fmt.Printf("Pass    : %v\n", pow.Validate())
	}
}
