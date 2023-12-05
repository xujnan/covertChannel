package btc

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func GenerateKey() {
	// 选择比特币网络，比如主网
	netParams := &chaincfg.MainNetParams

	// 创建一个新的随机私钥
	//NewPrivateKey instantiates a new private key from a scalar encoded as a big integer
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// 从私钥获取对应的公钥
	pubKey := privKey.PubKey()

	// 获取比特币地址
	address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), netParams)
	if err != nil {
		log.Fatalf("Failed to create address: %v", err)
	}

	fmt.Printf("Private Key: %x\n", privKey)
	fmt.Printf("Public Key: %x\n", pubKey.SerializeCompressed())
	fmt.Printf("Bitcoin Address: %s\n", address.EncodeAddress())

	// 你还可以导出私钥的WIF格式
	wif, err := btcutil.NewWIF(privKey, netParams, true)
	if err != nil {
		log.Fatalf("Failed to create WIF: %v", err)
	}
	fmt.Printf("WIF: %s\n", wif.String())

	// 你可以使用私钥签名交易等操作
	// 这里只是生成地址、公钥和私钥的示例
}
