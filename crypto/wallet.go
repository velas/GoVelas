package crypto

import (
	"bytes"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Wallet object, contain address as byte array and in base58 encoding
type Wallet struct {
	Base58Address string
	Address       []byte
}

// Version of HD wallet(VL)
var Version = []byte{15, 244}

// Address checksum length
const addressChecksumLen = 4

// Create wallet from public key
func CreateWallet(pubKey []byte) (*Wallet, error) {
	publickHash, err := hashPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	versionedPayload := append(Version, publickHash...)
	checksum := checksum(versionedPayload)

	address := append(versionedPayload, checksum...)

	return &Wallet{
		Base58Address: base58.Encode(address),
		Address:       address,
	}, nil
}

// hash of public key for generate wallet
func hashPublicKey(pubKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil, err
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160, nil
}

// checksum of hash
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// Check string is wallet
func IsWalletAddress(addr string) bool {
	address := base58.Decode(addr)

	lenAddress := len(address)
	payload := address[:lenAddress-addressChecksumLen]
	version := address[:2]
	if !bytes.Equal(version, Version) {
		return false
	}
	check := checksum(payload)
	checkAddress := append(payload, check...)
	if !bytes.Equal(address, checkAddress) {
		return false
	}
	return true
}
