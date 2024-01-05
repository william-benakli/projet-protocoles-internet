package cryptographie

/*
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

Pour obtenir la clé publique associée :
func test2() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	//Pour formater la clé publique comme une chaîne de 64 octets :
	formatted := make([]byte, 64)
	publicKey.X.FillBytes(formatted[:32])
	publicKey.Y.FillBytes(formatted[32:])
	//Pour parser une clé publique représentée comme une chaîne de 64 octets :
	var x, y big.Int
	x.SetBytes(data[:32])
	y.SetBytes(data[32:])
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
	hashed := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed[:])
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	//Pour vérifier un message :
	var r, s big.Int
	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:])
	hashed := sha256.Sum256(data)
	ok = ecdsa.Verify(publicKey, hashed[:], &r, &s)

}
*/
