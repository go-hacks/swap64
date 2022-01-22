// Everything key related

package main

import (
	"encoding/binary"
	//"fmt"
)

const maxVal64 uint64 = 18446744073709551615
// 0-63 base vals
var slotVals = []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63}
// Every other prime number (nothing up sleeve)
var pkSlots = []int{3,7,13,19,29,37,43,53}

// mashes seed with initial key
func seedKey (key []byte, seed []byte) []byte {
	var cnt int = 0
	for i := 0; i < len(key); i++ {
		if cnt == len(seed) {
			cnt = 0
		}
		key[i] = key[i] ^ seed[cnt]
		cnt++
	}
	return key
}

// Key schedule uses an 8 gate NLFSR.
// XOR + addition, then bitflip(which is still xor).
// Gate positions change every round based on key
// values in every other prime number slot.
func keyMaker (key []uint64) []uint64 {
	slots := make([]int, 64)
	copy(slots,slotVals)
	gates := make([]int, 8)
	for i := 0; i < 8; i++ {
		gatePos := key[pkSlots[i]] % uint64(len(slots))
		gates[i] = slots[int(gatePos)]
		slots = spliceInt(slots,int(gatePos))
	}
	xorPairs := make([]uint64, 4)
	for i := 0; i < blockSize64; i++ {
		xorPairs[0] = key[gates[0]] ^ key[gates[7]]
		xorPairs[1] = key[gates[1]] ^ key[gates[6]]
		xorPairs[2] = key[gates[2]] ^ key[gates[5]]
		xorPairs[3] = key[gates[3]] ^ key[gates[4]]
		new64 := xorPairs[0] + xorPairs[1] + xorPairs[2] + xorPairs[3]
		key = append(key[1:64], key[0:1]...)
		key[63] = new64 ^ maxVal64
	}
	return key
}

func keyExpansion (pass []byte) []byte {
	key := make([]byte, 512)
	key = pass
	sendByte := make([]byte, 1)
	pos := 0
	for {
		sendByte[0] = key[pos]
		expBytes := sbox8To16Bit(sendByte)
		key = append(key, expBytes...)
		if len(key) >= len(pass)+512 {
			key = key[len(key)-512:len(key)]
			break
		}
		pos++
	}
	return key
}

// hash keys for separate threads and initialize original key
func hasher (key []byte) []byte {
	sboxArr1, sboxArr2 := genSboxArrays(key)
	key = keyCubeMod(key)
	key64 := byteToUint64(key)
	for i := 0; i < 32; i++ {
		key64 = keyMaker(key64)
	}
	key = uint64ToByte(key64)
	key = sbox8To8Bit(key,sboxArr1,sboxArr2)
	key = flipBits8(key)
	key64 = byteToUint64(key)
	for i := 0; i < 32; i++ {
		key64 = keyMaker(key64)
	}
	key = uint64ToByte(key64)
	return key
}

// Cubes each byte of key and mods into range (one way function).
// Note: If over used, creates too many zeroes in key.
func keyCubeMod (arr []byte) []byte {
	for i, v := range arr {
		arr[i] = byte(int(v * v * v) % 256 )
	}
	return arr
}

// generates the 8bit sbox arrays
func genSboxArrays (key []byte) ([256]byte, [256]byte){
	baseArr1 := make([]int, 256)
	baseArr2 := make([]int, 256)
	for i := 0; i < 256; i++ {
		baseArr1[i] = i
		baseArr2[i] = i
	}
	var arr1 [256]byte//:= make([]byte, 256)
	var arr2 [256]byte// := make([]byte, 256)
	a := 511
	for i := 0; i < 256; i++ {
		// These should be faster than converting
		// to float64s and importing math.Pow()
		val1 := int(key[i]) * int(key[i]) * int(key[i]) % len(baseArr1)
		val2 := int(key[a]) * int(key[a]) * int(key[a]) % len(baseArr2)
		arr1[i] = byte(baseArr1[val1])
		arr2[i] = byte(baseArr2[val2])
		baseArr1 = spliceInt(baseArr1, val1)
		baseArr2 = spliceInt(baseArr2, val2)
		a--
	}
	return arr1, arr2
}

func sbox8To8Bit (key []byte, arr1, arr2 [256]byte) []byte {
	for i := 0; i < len(key); i++ {
		for j := 0; j < len(arr1); j++ {
			if key[i] == arr1[j] {
				key[i] = arr2[j]
				break
			}
		}
	}
	return key
}

func sbox8To16Bit (inData []byte) []byte {
	var sqrt3Hex string = "1c98c677de371c7dce08f7b716f08f566d94d23294ae1f771f851f094b82420c1f4090d9391a1a656119ddb4661da2441a2f0bb426d9251b22f8bb0940ab2f3e91f42150246b4a4f15ce1a8dbc49043c1a5d5ae105fb5b881a2281cc21d91a42d58c2402945358b01bdb19cb1a823e084a6926899203219122f21de287be19b013d319da87fc2586276e1db2491f0cb2207ae3b75d97208926322660c68c220ab52f71f5b878400eb336d8651a17e2cf1ff7d1a22b4c7bf91f99ab073ee6227ed7d1068609301bc7ad001b2f2387f8a925ec1c00c28f71422471fd4a206e1eb8219b269f9729d1fe5a8545a79aaf74861e9e1b741e1998dd1b2c81610e7b6721242129b420274dd820b7fa72378c41661d61094426d61bab2966267525032694fb202311258e1e9d256f233e229621510df80d396db0fb7d25fbae75175025fa2218fef03c166816267200d3254ca2501ea42e2421cdb00e8977fde91a217301d618b3f76c429c3df5ce209a7389ab3c4a07278f3a0f222a2625b39eb07a216626191edd9796fdc1239ff30e2563293b5ee4916f24031e3a66903fba1ae44f68d1d41aad96f9ef071b891cdc8e2bc1b91fd11d981c1183c88165239426edd3e661021ebf2620cc87263023e87e65089b14b620d3218fd3851bc35ee7e7461b9609cdcb31e8ef218b6c0fea476ece86061f4f1a02d1cefcdc215af480584d201e"
	var sbox [256]uint16
	for i := 0; i < 256; i ++ {
		str, _ := substr(sqrt3Hex,i*4,4)
		sbox[i] = hexToUint16([]byte(str))
	}
	outData := make([]byte, len(inData)*2)
	for i := 0; i < len(inData); i++ {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes,sbox[int(inData[i])])
		outData[i*2] = bytes[0]
		outData[i*2+1] = bytes[1]
	}
	return outData
}
