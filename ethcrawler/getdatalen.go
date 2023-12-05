package ethcrawler

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 获取交易中data字段的长度，单位是字节
func GetTxDataLen() {
	// 以太坊节点的URL
	ethereumURL := "https://eth-mainnet.g.alchemy.com/v2/0CN38V9wrCa7IQ4I86wy0beAn7C8lQHk"

	// 连接到以太坊节点
	client, err := ethclient.Dial(ethereumURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum ethclient: %v", err)
	}

	//number, err := ethclient.BlockNumber(context.Background())
	//if err != nil {
	//	return
	//}
	number := 18484392
	for j := 0; j < 8; j++ {

		var blocks []*types.Block
		for i := number; i > number-10; i-- {
			// 获取指定区块的信息（替换块号）
			blockNumber := big.NewInt(int64(i))

			block, err := client.BlockByNumber(context.Background(), blockNumber)
			if err != nil {
				log.Fatalf("Failed to retrieve block %s: %v", blockNumber.String(), err)
			}

			fmt.Println(block.Number())
			blocks = append(blocks, block)
		}

		fileName := fmt.Sprint(number-9) + "-" + fmt.Sprint(number)
		// 打开文件以追加数据
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close() // 确保在函数结束时关闭文件

		for _, b := range blocks {
			txs := b.Transactions()
			for _, tx := range txs {
				dataLen := len(tx.Data())
				// 将数据追加到文件
				_, err = file.WriteString(strconv.Itoa(dataLen) + "\n")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		number -= 10
	}

}
