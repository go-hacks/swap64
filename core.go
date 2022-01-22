// Core cipher functions

package main

import (
	"sync"
	//"fmt"
)

// main process (fwd and rev)
func algo (sd sectionData, key []uint64, threadNum int, wg *sync.WaitGroup) {
	partBlockCnt := sd.blockCnt / sd.threads
	chunkByteSize := partBlockCnt * blockSize
	var numBlocks int
	if sd.blockCnt == 1 {
		numBlocks = 1
	} else {
		numBlocks = sd.blockCnt / sd.threads
		if threadNum == sd.threads {
			numBlocks = sd.blockCnt - ((sd.threads - 1) * partBlockCnt)
		}
	}
	var isLastBlock bool = false
	if sd.threads == threadNum && sd.remain > 0 && opMode == "fwd" {
		isLastBlock = true
	}
	dataBytes := make([]byte, blockSize)
	dataBlock := make([]uint64, blockSize64)
	bufBlock := make([]uint64, blockSize64)
	cbcBlock := make([]uint64, blockSize64)
	for i := 1; i <= numBlocks; i++ {
		filePos := int64(((threadNum - 1) * chunkByteSize) + ((i - 1) * blockSize))
		if i == numBlocks && isLastBlock {
			lastBytes := make([]byte, sd.remain)
			_, _ = sd.inFile.ReadAt(lastBytes, filePos)
			for i := 0; i < len(lastBytes); i++ {
				dataBytes[i] = lastBytes[i]
			}
			fillBytes := getRandData(blockSize-sd.remain)
			dataBytes = append(dataBytes[0:sd.remain], fillBytes...)
		} else {
			_, _ = sd.inFile.ReadAt(dataBytes, filePos)
		}
		dataBlock = byteToUint64(dataBytes)
		switch opMode {
			case "fwd":
				xor64(dataBlock, cbcBlock)
				dataBlock, key = cipherFwd(dataBlock, key)
				copy(cbcBlock,dataBlock)
			case "rev":
				copy(bufBlock,dataBlock)
				dataBlock, key = cipherRev(dataBlock, key)
				xor64(dataBlock, cbcBlock)
				copy(cbcBlock,bufBlock)
		}
		dataBytes = uint64ToByte(dataBlock)
		if sd.threads == threadNum && i == numBlocks && opMode == "rev" {
			dataBytes = dataBytes[0:blockSize-sd.trimVal]
		}
		_, _ = sd.outFile.WriteAt(dataBytes, filePos)
	}
	wg.Done()
	return
}

// encryption process for each block
func cipherFwd (dataBlock []uint64, key []uint64) ([]uint64, []uint64) {
	keyAry := makeMatrix(64, roundCnt)
	keyAry[0] = keyMaker(key)
	for i := 1; i < roundCnt; i++ {
		keyAry[i] = keyMaker(keyAry[i-1])
	}
	swapVals := swapValGen(keyAry[0])
	dataBlock = swap64Fwd(dataBlock, swapVals)
	for i := 0; i < roundCnt; i++ {
		xor64(dataBlock, keyAry[i])
		dataBlock = append(dataBlock[8:blockSize64], dataBlock[0:8]...)
	}
	return dataBlock, keyAry[len(keyAry)-1]
}

// decryption process for each block
func cipherRev (dataBlock []uint64, key []uint64) ([]uint64, []uint64) {
	keyAry := makeMatrix(64, roundCnt)
	keyAry[0] = keyMaker(key)
	for i := 1; i < roundCnt; i++ {
		keyAry[i] = keyMaker(keyAry[i-1])
	}
	for i := roundCnt-1; i >= 0; i-- {
		dataBlock = append(dataBlock[blockSize64-8:blockSize64], dataBlock[0:blockSize64-8]...)
		xor64(dataBlock, keyAry[i])
	}
	swapVals := swapValGen(keyAry[0])
	dataBlock = swap64Rev(dataBlock, swapVals)
	return dataBlock, keyAry[len(keyAry)-1]
}

func swap64Fwd(dataBlock []uint64, swapVals []int) []uint64 {
	newBlock := make([]uint64, blockSize64)
	for i := 0; i < blockSize64; i++ {
		newBlock[swapVals[i]] = dataBlock[i]
	}
	return newBlock
}

func swap64Rev(dataBlock []uint64, swapVals []int) []uint64 {
	newBlock := make([]uint64, blockSize64)
	for i := 0; i < blockSize64; i++ {
		newBlock[i] = dataBlock[swapVals[i]]
	}
	return newBlock
}

func swapValGen (key []uint64) ([]int) {
// hardcode this array
	baseVals := make([]int, 64)
	for i := 0; i < 64; i++ {
		baseVals[i] = i
	}
	vals := make([]int, 64)
	for i, v := range key {
		slot := int(v % uint64(len(baseVals)))
		vals[i] = baseVals[slot]
		baseVals = spliceInt(baseVals,slot)
	}
	return vals
}
