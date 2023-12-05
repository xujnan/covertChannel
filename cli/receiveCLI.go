package cli

import (
	"cc/ethclient"
	"cc/hashchain"
	"fmt"
)

func ReceiveCLI() {
	//初始化
	//设置共享公钥
	//生成密钥链、地址链
	pubKey := hashchain.GetPubKeyFromFile()
	hc := hashchain.NewHashChain()
	addChain := hashchain.NewAddressChain(pubKey) //通过共享公钥生成地址链
	ethClient := ethclient.NewEthClient()

	text := ethClient.ParseData(hc, addChain)
	fmt.Printf("\n数据：%s", text)
}
