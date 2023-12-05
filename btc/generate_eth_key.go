package btc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKeEthKey() {
	// 生成随机私钥
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// 获取公钥
	publicKey := privateKey.Public()

	// 通过公钥获取以太坊地址
	publicKeyBytes := crypto.FromECDSAPub(publicKey.(*ecdsa.PublicKey))
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	// 将私钥和地址打印出来
	privateKeyHex := hexutil.Encode(crypto.FromECDSA(privateKey))
	publicKeyHex := hexutil.Encode(publicKeyBytes)
	addressString := address.String()
	addressHex := address.Hex()

	fmt.Printf("私钥: %s\n", privateKeyHex)
	fmt.Printf("公钥: %s\n", publicKeyHex)
	fmt.Println(addressHex)
	fmt.Println(addressString)
}
