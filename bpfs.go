// Boiler plate functions

package main

import (
  "fmt"
  "encoding/binary"
  "encoding/hex"
  "strconv"
  "runtime"
  "os"
)

// Converts hex chars to raw byte values
func hexToByte(src []byte) []byte {
	dst := make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(dst, src)
	return dst
}

// Converts hex chars to integer
func hexToInt(hex_byte []byte) int {
	i, err := strconv.ParseInt(string(hex_byte), 16, 64)
	parseFatal(err, "")
	integer := int(i)
	return integer
}

// Converts 64bit array to byte array
func uint64ToByte(src []uint64) []byte {
	byteSize := len(src) * 8
	dst := make([]byte, byteSize)
	for i := 0; i < len(src); i++ {
		binary.LittleEndian.PutUint64(dst[i*8:i*8+8], src[i])
	}
	return dst
}

// Converts byte array to 64bit array
func byteToUint64(src []byte) []uint64 {
	if len(src)%8 != 0 {
		fmt.Println("Incorrect length of input to ByteToUint64!")
		os.Exit(1)
	}
	len64 := len(src) / 8
	dst := make([]uint64, len64)
	for i := 0; i < len64; i++ {
		dst[i] = binary.LittleEndian.Uint64(src[i*8:i*8+8])
	}
	return dst
}

// Perl like substring functionality
func substr(s string, start int, str_len int) (string, string) {
	var err string
	if start+str_len > len(s) {
		err = "ERROR:Substring length > length of string"
	}
	end := start + str_len
	if str_len < 0 {
		end = len(s) + str_len
	}
	if start < 0 {
		end = len(s) - (len(s) + start)
		start = 0
	}
	arr := []byte(s)
	slice := arr[start:end]
	return string(slice), err
}

// Fatal error parser
func parseFatal(err error, msg string) {
	if err != nil {
		if msg != "" {
			fmt.Println(msg)
		}
		fmt.Println(err)
		os.Exit(0)
	} else {
		return
	}
}

// Gets size of file fname
func getFileSize(fname string) int64 {
	fi, err := os.Stat(fname)
	parseFatal(err, "")
	fileSize := fi.Size()
	return fileSize
}

// Returns number of CPUs in machine
func threads() int {
	return runtime.NumCPU()
}

// Makes a 64bit matrix of given dimensions
func makeMatrix (x int, y int) [][]uint64 {
	matrix := make([][]uint64, y)
	for i := range matrix {
		matrix[i] = make([]uint64, x)
	}
	return matrix
}

// Removes given index from byte array
func splice(arr []byte, index int) []byte {
	return append(arr[0:index], arr[index+1:len(arr)]...)
}

// Removes given index from integer array
func spliceInt(arr []int, index int) []int {
	return append(arr[0:index], arr[index+1:len(arr)]...)
}

// XORs two byte arrays
func xor(arr1 []byte, arr2 []byte) []byte {
	for i, v := range arr1 {
		arr1[i] = v ^ arr2[i]
	}
	return arr1
}

// XORs two 64bit arrays
func xor64(arr1 []uint64, arr2 []uint64) {
	for i, v := range arr1 {
		arr1[i] = v ^ arr2[i]
	}
	return
}

// Flips bits of byte array by XORing with all 1s
func flipBits8 (arr []byte) []byte {
  for i, _ := range arr {
    arr[i] ^= 255
  }
  return arr
}

// Converts 4hex chars into uint16
func hexToUint16(hexByte []byte) uint16 {
	i, _ := strconv.ParseInt(string(hexByte), 16, 64)
	integer := uint16(i)
	return integer
}
