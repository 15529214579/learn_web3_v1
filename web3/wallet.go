package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	//PubKey *ecdsa.PublicKey
	//为了实现方便，这里的publickey存储rs串
	PubKey []byte
}

func NewWallet() *Wallet {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}

	pubKeyOrig := privateKey.PublicKey
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)
	return &Wallet{Private: privateKey, PubKey: pubKey}
}

func (w *Wallet) NewAddress() string {
	pubkey := w.PubKey
	hash := HashPubKey(pubkey)
	version := byte(00)

	payload := append([]byte{version}, hash...)

	checkCode := CheckSum(payload)

	payload = append(payload, checkCode...)

	address := base58.Encode(payload)
	return address
}

func HashPubKey(data []byte) []byte {
	//此处代码copy,后续可替换成其他hash算法
	hash := sha256.Sum256(data)

	return hash[:]
}

func CheckSum(data []byte) []byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])

	//返回前4位校验码
	checkCode := hash2[:4]
	return checkCode
}

func IsValidAddress(address string) bool {
	//1. 解码
	addressByte := base58.Decode(address)

	if len(addressByte) < 4 {
		return false
	}

	//2. 取数据
	payload := addressByte[:len(addressByte)-4]
	checksum1 := addressByte[len(addressByte)-4:]

	//3. 做checksum函数
	checksum2 := CheckSum(payload)

	fmt.Printf("checksum1 : %x\n", checksum1)
	fmt.Printf("checksum2 : %x\n", checksum2)

	//4. 比较
	return bytes.Equal(checksum1, checksum2)
}
