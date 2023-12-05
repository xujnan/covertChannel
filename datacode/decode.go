package datacode

import (
	"log"
	"strconv"
)

func GetDataLens(data [][]byte) []int {
	var lens []int
	for _, d := range data {
		lens = append(lens, len(d))
	}

	return lens
}

func LensToCipehrBytes(table *EncodeTable, lens []int) []byte {
	binaryString := ""
	for _, l := range lens {
		binaryString += table.LTBTable[l]
	}

	// 初始化一个字节数组
	byteData := make([]byte, len(binaryString)/8)

	// 将二进制字符串解析为字节数组
	for i := 0; i < len(byteData); i++ {
		byteStr := binaryString[i*8 : (i+1)*8]
		byteValue, err := strconv.ParseUint(byteStr, 2, 8) //二进制有个符号位这个概念
		if err != nil {
			log.Panic("Error parsing binary string:", err)
		}

		byteData[i] = byte(byteValue)
	}

	return byteData
}
