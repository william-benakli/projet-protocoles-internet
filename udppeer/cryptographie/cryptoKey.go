package cryptographie

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

var PrivateKey *ecdsa.PrivateKey
var PublicKey *ecdsa.PublicKey
var OtherPublicKey ecdsa.PublicKey // clef de l'interlocuteur

func GeneratePrivateKey() *ecdsa.PrivateKey {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return privateKey
}
func GetPublicKey(privateKey *ecdsa.PrivateKey) *ecdsa.PublicKey {
	publicKey, _ := privateKey.Public().(*ecdsa.PublicKey)
	return publicKey
}

func FormateKey() []byte {
	formatted := make([]byte, 64)
	PublicKey.X.FillBytes(formatted[:32])
	PublicKey.Y.FillBytes(formatted[32:])
	return formatted
}
func UnFormateKey(data []byte) ecdsa.PublicKey {
	var x, y big.Int
	x.SetBytes(data[:32])
	y.SetBytes(data[32:])
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
	return publicKey
}

func Encrypted(data []byte) []byte {
	hashed := sha256.Sum256(data)
	r, s, _ := ecdsa.Sign(rand.Reader, PrivateKey, hashed[:])
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	return signature
}

func FormatePrivateKey() []byte {
	formatted := make([]byte, 64)
	PrivateKey.X.FillBytes(formatted[:32])
	PrivateKey.Y.FillBytes(formatted[32:])
	return formatted
}

func VerifyHash(data []byte, signature []byte) bool {
	var r, s big.Int
	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:])
	hashed := sha256.Sum256(data)
	ok := ecdsa.Verify(&OtherPublicKey, hashed[:], &r, &s)
	return ok
}
