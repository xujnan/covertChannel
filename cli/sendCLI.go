package cli

import (
	"bufio"
	"cc/ethclient"
	"cc/hashchain"
	"fmt"
	"os"
)

func SendCLI() {
	pubKey := hashchain.GetPubKeyFromFile()
	hashChain := hashchain.NewHashChain()
	addressChain := hashchain.NewAddressChain(pubKey)
	ethClient := ethclient.NewEthClient()

	//输入传输数据
	fmt.Print("输入数据：")
	reader := bufio.NewReader(os.Stdin)
	data, _ := reader.ReadString('\n')

	//TODO 密钥处理
	ethClient.SendTxsWithData(data, hashChain, addressChain)
}
