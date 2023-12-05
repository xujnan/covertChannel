package datacode

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const NLEN = 3 //表示选择八个长度，一个长度代表3个二进制数
const LENGTHS = 8

type EncodeTable struct {
	Table    map[string]int
	LTBTable map[int]string
}

// NewEncodeTable len(lens) = 2的次方
func NewEncodeTable() *EncodeTable {
	table := make(map[string]int)
	LTBTable := make(map[int]string)
	lens := GetHighFrequencyLens()

	for i, l := range lens {
		bi := fmt.Sprintf("%03b", i)
		table[bi] = l
		LTBTable[l] = bi
	}

	return &EncodeTable{table, LTBTable}
}

func BytesToBinaryString(data []byte) string {
	// 输出二进制字符串
	binaryString := ""
	for _, b := range data {
		binaryString += fmt.Sprintf("%08s", strconv.FormatUint(uint64(b), 2))
	}

	return binaryString
}

// 每个分片得到对应的长度
func BinaryDataToLens(table *EncodeTable, binaryData string) []int {
	split := Split(binaryData)
	var lens []int

	for _, v := range split {
		l := table.Table[v]
		lens = append(lens, l)
	}

	return lens
}

// 获得对应长度的随机数据
func LensToRandomData(lens []int) [][]byte {
	var randomBytes [][]byte

	for _, length := range lens {
		bytes := make([]byte, length)
		_, err := rand.Read(bytes)
		if err != nil {
			fmt.Println("随机数生成失败:", err)
			return nil
		}
		randomBytes = append(randomBytes, bytes)
	}

	return randomBytes
}

// 将二进制流分片
func Split(binaryData string) []string {
	nLen := NLEN
	var splitData []string

	if len(binaryData)%nLen != 0 {
		binaryData = binaryData + strings.Repeat("0", nLen-(len(binaryData)%nLen))
	}

	for i := 0; i < len(binaryData); i += nLen {
		str := binaryData[i : i+nLen]
		splitData = append(splitData, str)
	}

	return splitData
}

// 从文件中获得最高频长度
func GetHighFrequencyLens() []int {
	var lens []int

	file, err := os.Open("./data/mr-out-0")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return nil
	}
	defer file.Close() // 在函数结束时关闭文件

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		line := scanner.Text() // 从扫描器中读取一行
		l := strings.Split(line, " ")[0]
		ll, _ := strconv.Atoi(l)
		lens = append(lens, ll)

		i++
		if i == LENGTHS {
			break
		}
	}

	return lens
}
