package ethclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// 定义结构体以映射 JSON 数据
type Transfer struct {
	BlockNum string `json:"blockNum"`
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type Result struct {
	Transfers []Transfer `json:"transfers"`
}

type Response struct {
	Result Result `json:"result"`
}

// 通过地址获取交易hash
// 通过调用浏览器API实现
func (c *Client) GetTxHashesByAddr(addresses []string) [][]byte {
	var txHashes [][]byte

	for i, addr := range addresses {
		payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"alchemy_getAssetTransfers\",\"params\":[{\"fromBlock\":\"0x0\",\"toBlock\":\"latest\"," + "\"toAddress\":\"" +
			addr +
			"\",\"withMetadata\":false,\"excludeZeroValue\":true,\"maxCount\":\"0x3e8\"," +
			"\"category\":[\"external\"],\"order\":\"desc\"}]}")
		req, _ := http.NewRequest("POST", URL, payload)
		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Panic(err)
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		var r Response
		// 解析 JSON 字符串
		err = json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Panic("解析 JSON 字符串时出错:", err)
		}

		txHash := r.Result.Transfers[0].Hash
		txHashes = append(txHashes, []byte(txHash))

		fmt.Print(i+1, " ")
	}

	return txHashes
}
