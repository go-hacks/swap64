// Functions for adding extra data at eof

package main

import (
	"fmt"
	"os"
)

// write encrypted seed to eof
func addSeed (outFile *os.File, seed []byte) {
	outFile.Seek(0, 2)
	outFile.Write(seed)
}

// write fill count to eof
func addHexFill (outFile *os.File, remain int) {
	var hexFill string
	if remain > 0 {
		hexFill = fmt.Sprintf("%04x", blockSize - remain)
	} else {
		hexFill = "0000"
	}
	fillCnt := hexToByte([]byte(hexFill))
	outFile.Seek(0, 2)
	outFile.Write(fillCnt)
}

//write thread cnt to eof
func addThreadCnt (outFile *os.File, threads int) {
	outFile.Seek(0, 2)
	outByte := make([]byte,1)
	outByte[0] = byte(uint8(threads))
	outFile.Write(outByte)
	return
}
