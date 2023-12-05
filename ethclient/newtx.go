package ethclient

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	EthClient *ethclient.Client
}

type Transaction struct {
	Data  []byte
	Value *big.Int
	To    string
}

// URL https://eth-mainnet.g.alchemy.com/v2/0CN38V9wrCa7IQ4I86wy0beAn7C8lQHk
// HTTP://127.0.0.1:7545
// https://eth-sepolia.g.alchemy.com/v2/Bf4Vzi4zCrPqqgNtckE1ygzcY0e9oLnm
const URL = "https://eth-sepolia.g.alchemy.com/v2/Bf4Vzi4zCrPqqgNtckE1ygzcY0e9oLnm"

func NewTransaction(data []byte, value big.Int, to string) *Transaction {
	//value.Mul(&value, big.NewInt(1000000000000000000))
	tx := &Transaction{
		Data:  data,
		Value: &value,
		To:    to,
	}

	return tx
}

func NewEthClient() *Client {
	client, err := ethclient.Dial(URL)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{client}
}

func (c *Client) SendTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := tx.Value                 // in wei (0.01 ethcrawler)
	gasLimit := uint64(70000)         // in units
	tip := big.NewInt(1000000000)     // maxPriorityFeePerGas = 1 Gwei
	feeCap := big.NewInt(20000000000) // maxFeePerGas = 20 Gwei  总费用 = (燃气用量) x （(燃气单价) + Gas Tip）
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(tx.To)

	chainID, err := c.EthClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	transaction := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      tx.Data,
	})

	signedTx, err := types.SignTx(transaction, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = c.EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("Transaction hash: %s", signedTx.Hash().String())
}
