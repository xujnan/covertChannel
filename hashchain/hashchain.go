package hashchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

type HashChain struct {
	K           *big.Int //需要共享的密钥K
	Chain       [][]byte //K0 -- Kn
	Length      int
	LastEncrypt int //最新可用密钥索引
	LastDecrypt int //最新可用密钥索引
	LastAddr    int
}

type AddressChain struct {
	SharedPubKey *ecdsa.PublicKey   //共享公钥
	PubKeyChain  []*ecdsa.PublicKey //不包括共享密钥
	LastAddr     int                //最后未使用地址索引
}

type SharedPubPriKey struct {
	SharedSK *ecdsa.PrivateKey
}

// NewHashChain 新建一条密钥链
func NewHashChain() *HashChain {
	//privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	//if err != nil {
	//	log.Fatal(err)
	//}
	f, _ := os.Open("./cli/key")
	privateKey, err := io.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	key, err := crypto.HexToECDSA(fmt.Sprintf("%x", privateKey))
	if err != nil {
		log.Panic(err)
	}

	k := key.D
	chain := [][]byte{k.Bytes()}
	hashChain := HashChain{k, chain, 1, 1, 1, 1}

	return &hashChain
}

func (hc *HashChain) ExtendHashChain(l int) {
	lastIndex := len(hc.Chain) - 1

	for i := 0; i < l; i++ {
		k := sha256.Sum256(hc.Chain[lastIndex])
		newK := mod(k[:])
		hc.Chain = append(hc.Chain, newK[:])
		lastIndex++
	}

	hc.Length += l
}

func (hc *HashChain) GetEncryptKey() []byte {
	if hc.Length-hc.LastEncrypt <= 0 {
		hc.ExtendHashChain(1)
	}

	key := hc.Chain[hc.LastEncrypt]
	hc.LastEncrypt++

	return key
}

func (hc *HashChain) GetDecryptKey() []byte {
	if hc.Length-hc.LastDecrypt <= 0 {
		hc.ExtendHashChain(1)
	}

	key := hc.Chain[hc.LastDecrypt]
	hc.LastDecrypt++

	return key
}

func NewAddressChain(publicKey *ecdsa.PublicKey) *AddressChain {
	addressChain := &AddressChain{
		SharedPubKey: publicKey,
		PubKeyChain:  []*ecdsa.PublicKey{},
		LastAddr:     0,
	}

	return addressChain
}

// ExtendPubKeyChain 从地址链上获得num个地址
func (ac *AddressChain) ExtendPubKeyChain(num int, hashChain *HashChain) {
	hashChain.ExtendHashChain(num)
	var publicKeys []*ecdsa.PublicKey

	for i := 0; i < num; i++ {
		x, y := crypto.S256().ScalarBaseMult(hashChain.Chain[hashChain.LastAddr]) //K*G
		Q := &ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}                   //曲线上新的一点，K*G的结果

		// 执行点加法操作生成新公钥 C
		xC, yC := crypto.S256().Add(ac.SharedPubKey.X, ac.SharedPubKey.Y, Q.X, Q.Y)
		publicKey := &ecdsa.PublicKey{Curve: crypto.S256(), X: xC, Y: yC}
		publicKeys = append(publicKeys, publicKey)
		hashChain.LastAddr++
	}
	ac.PubKeyChain = append(ac.PubKeyChain, publicKeys...)
}

// GetAddress 从地址链上获得num个地址
func (ac *AddressChain) GetAddress(num int, hashChain *HashChain) []string {
	var addresses []string
	ac.ExtendPubKeyChain(num, hashChain)

	for i := 0; i < num; i++ {
		address := crypto.PubkeyToAddress(*ac.PubKeyChain[ac.LastAddr])
		ac.LastAddr++
		addresses = append(addresses, address.Hex())
	}

	return addresses
}

func NewSharedPubPriKey() *SharedPubPriKey {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	sharedKey := &SharedPubPriKey{privateKey}

	return sharedKey
}

// 求两个大数模,使每个密钥k均满足椭圆曲线
func mod(k []byte) []byte {
	order := crypto.S256().Params().P
	m := new(big.Int)
	kBigInt := new(big.Int).SetBytes(k)

	_, m = new(big.Int).DivMod(kBigInt, order, m)

	return m.Bytes()
}
