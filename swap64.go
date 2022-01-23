// Multi-threaded, scalable, x64 CBC 4096-bit block
// cipher that performs 7 rounds per block with an
// NLFSR key schedule, key dependant movements, and
// a 512-bit OTP key seed masked w/passphrase hash.
// Version 0.77

package main

import (
	"fmt"
	"os"
	"sync"
)

// Declare global vars
const blockSize int = 512
const blockSize64 int = 64
var opMode string
var roundCnt int

type sectionData struct {
	inFile		*os.File
	outFile		*os.File
	threads		int
	blockCnt	int
	trimVal		int
	remain		int
}

func main () {
	// initialize startup values
	baseFileName := os.Args[len(os.Args)-1]
	var outFileName, opMode string
	var secData sectionData
	outFileName, opMode, secData.threads, roundCnt = parseParams()
	var fSize int64 = getFileSize(baseFileName)
	secData.remain = int(fSize) % blockSize
	// get password and generate key
	passphrase := getPassphrase()
	key := keyExpansion([]byte(passphrase))
	key = hasher(key)
	// open input file
	var inErr error
	secData.inFile, inErr = os.OpenFile(baseFileName,os.O_RDONLY,0666)
	parseFatal(inErr,"")
	// create passphrase/key based mask for random seed
	seedMask := make([]byte, 64)
	copy(seedMask,key[0:64])
	expMask := hasher(keyExpansion(seedMask))
	seedMask = expMask[448:blockSize]
	var seed []byte
	switch opMode {
		case "fwd":
			// generate random seed
			seed = getRandData(blockSize64)
		case "rev":
			// extract file information (encrypted seed,filler count,threads)
			fileData := make([]byte, 67)
			// Seek to before needed bytes at EOF.
			// No need to seek back after as all further
			// reads and writes happen at specific positions.
			secData.inFile.Seek(-67, 2)
			n, readErr := secData.inFile.Read(fileData)
			readErrMsg := fmt.Sprintf("read %d bytes & error", n)
			parseFatal(readErr, readErrMsg)
			encSeed := fileData[0:blockSize64]
			seed = xor(encSeed, seedMask)
			trimCnt := fileData[blockSize64:blockSize64+2]
			hexTrim := fmt.Sprintf("%02x",trimCnt)
			secData.trimVal = hexToInt([]byte(hexTrim))
			secData.threads = int(fileData[blockSize64+2])
	}
	// mask key with random seed
	// This makes the output file completely diff
	// even if it's the same file and passphrase.
	key = seedKey(key, seed)
	// open output file
	var outErr error
	secData.outFile, outErr = os.OpenFile(outFileName,os.O_CREATE|os.O_WRONLY,0666)
	parseFatal(outErr,"")
	// calculate blocks
	secData.blockCnt = int(fSize) / blockSize
	if secData.remain > 0 && opMode == "fwd" {
		secData.blockCnt++
	}
	if secData.blockCnt < secData.threads {
		secData.threads = 1
	}
	// generate keys for multithread
	keys := make([][]byte, secData.threads)
	keys64 := make([][]uint64, secData.threads)
	keys[0] = key
	keys64[0] = byteToUint64(key)
	for i := 1; i < secData.threads; i++ {
		keys[i] = hasher(keys[i-1])
		keys64[i] = byteToUint64(keys[i])
	}

	// main phase->encrypt or decrypt
	switch opMode {
		case "fwd":
			fmt.Printf("Encrypting via %d thread(s) w/ %d round(s) per block...", secData.threads, roundCnt)
			var wg sync.WaitGroup
			for i := 0; i < secData.threads; i++ {
				wg.Add(1)
				go algo(secData, keys64[i], i+1, &wg)
			}
			wg.Wait()
			// mask seed and store at eof
			encSeed := xor(seed,seedMask)
			addSeed(secData.outFile, encSeed)
			// add filler data count at eof
			addHexFill(secData.outFile, secData.remain)
			// add thread cnt at eof
			addThreadCnt(secData.outFile, secData.threads)
			fmt.Printf("Done.\n")
			fmt.Println("Encryption Complete!")
		case "rev":
			fmt.Printf("Decrypting via %d thread(s) w/ %d round(s) per block...", secData.threads, roundCnt)
			var wg sync.WaitGroup
			for i := 0; i < secData.threads; i++ {
				wg.Add(1)
				go algo(secData, keys64[i], i+1, &wg)
			}
			wg.Wait()
			fmt.Printf("Done.\n")
			fmt.Println("Decryption Complete!")
	}
	// end phase->close up shop
	outErr = secData.outFile.Close()
	parseFatal(outErr,"")
	inErr = secData.inFile.Close()
	parseFatal(inErr,"")
}
