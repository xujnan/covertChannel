package ethclient

import (
	"cc/AES"
	"cc/datacode"
	"cc/hashchain"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

// 获取所有交易里的data字段
func (c *Client) GetTxDatas(txHashes [][]byte) [][]byte {
	var data [][]byte

	for _, hash := range txHashes {
		transactionHash := common.HexToHash(string(hash))

		tx, _, err := c.EthClient.TransactionByHash(context.Background(), transactionHash)
		if err != nil {
			log.Panic(err)
		}

		data = append(data, tx.Data())
	}

	return data
}

// 解析第一笔交易。获取后续交易数量
func (c *Client) ParseData(hc *hashchain.HashChain, addChain *hashchain.AddressChain) string {
	addr := addChain.GetAddress(1, hc)
	firshTxHash := c.GetTxHashesByAddr(addr)
	firstData := c.GetTxDatas(firshTxHash)
	txNum, err := strconv.Atoi(string(firstData[0]))
	if err != nil {
		log.Panic(err)
	}

	return c.ParseAllData(hc, addChain, txNum)
}

func (c *Client) ParseAllData(hc *hashchain.HashChain, addChain *hashchain.AddressChain, txNum int) string {
	// 解析后续交易
	addresses := addChain.GetAddress(txNum, hc)
	fmt.Printf("一共有%d笔交易\n", len(addresses))
	fmt.Print("查询并解析交易：")
	txHashes := c.GetTxHashesByAddr(addresses)
	txData := c.GetTxDatas(txHashes)

	table := datacode.NewEncodeTable()
	dataLens := datacode.GetDataLens(txData)
	cipher := datacode.LensToCipehrBytes(table, dataLens)

	decryptKey := hc.GetDecryptKey()
	plainText := AES.Decrypt(cipher, decryptKey)

	return string(plainText)
}
