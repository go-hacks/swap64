// Random data generation
package main

import (
	"bufio"
	"os"
)

// returns random data from /dev/random
// used for seed and filler data
func getRandData (count int) []byte {
	randFile, randErr := os.OpenFile("/dev/random",os.O_RDONLY,0666)
	parseFatal(randErr,"")
	randReader := bufio.NewReader(randFile)
	randBlob := make([]byte, count)
	_, readErr := randReader.Read(randBlob)
	parseFatal(readErr,"")
	randErr = randFile.Close()
	parseFatal(randErr,"")
	return randBlob
}
