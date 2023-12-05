package ethclient

import (
	"bufio"
	"cc/AES"
	"cc/datacode"
	"cc/hashchain"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

func SetTxsWithData(data string, hashChain *hashchain.HashChain, addressChain *hashchain.AddressChain) []*Transaction {
	table := datacode.NewEncodeTable()
	key := hashChain.GetEncryptKey()
	cipher := AES.Encrypt([]byte(data), key)
	binaryCipher := datacode.BytesToBinaryString(cipher)
	lens := datacode.BinaryDataToLens(table, binaryCipher)
	randomByte := datacode.LensToRandomData(lens)
	addresses := addressChain.GetAddress(len(randomByte)+1, hashChain)

	fmt.Printf("需要发送%d笔交易，输入每笔交易的发送金额：单位(wei)", len(randomByte)+1)
	reader := bufio.NewReader(os.Stdin)
	v, _ := reader.ReadString('\n')
	var value big.Int
	var txs []*Transaction
	value.SetString(v, 10)
	firstTx := NewTransaction([]byte(strconv.Itoa(len(randomByte))), value, addresses[0])
	txs = append(txs, firstTx)

	for i := 0; i < len(randomByte); i++ {
		tx := NewTransaction(randomByte[i], value, addresses[i+1])
		txs = append(txs, tx)
	}

	return txs
}

func (c *Client) SendTxsWithData(data string, hashChain *hashchain.HashChain, addressChain *hashchain.AddressChain) {
	txs := SetTxsWithData(data, hashChain, addressChain)

	privateKey, err := crypto.HexToECDSA("39a56a6bc3fbb7b5115b69cbcb8267f4dfc8fc2c073d7ac03096be36eb481005")
	if err != nil {
		log.Fatal(err)
	}

	for i, tx := range txs {
		fmt.Print(i+1, " ")
		c.SendTransaction(tx, privateKey)
	}
}
