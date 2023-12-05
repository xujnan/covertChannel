package hashchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func GetPubKeyFromFile() *ecdsa.PublicKey {
	f, _ := os.Open("./cli/publicKey")
	pub, err := io.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}
	x := new(big.Int).SetBytes(pub[:len(pub)/2])
	y := new(big.Int).SetBytes(pub[len(pub)/2:])

	pubKey := &ecdsa.PublicKey{}
	pubKey.Curve = crypto.S256()
	pubKey.X = x
	pubKey.Y = y

	return pubKey
}

func SetKToFile() {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./cli/key", privateKey.D.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func SetPubKeyToFile(publicKey *ecdsa.PublicKey) {
	var xy []byte

	xy = append(xy, publicKey.X.Bytes()...)
	xy = append(xy, publicKey.Y.Bytes()...)

	err := os.WriteFile("./cli/publicKey", xy, 0644)
	if err != nil {
		log.Panic(err)
	}
}
