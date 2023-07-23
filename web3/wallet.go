package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
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
