// Package crypto implements crypto algorithms for work with Velas
//
// Can create keys, wallet and transaction
package crypto

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"github.com/GoKillers/libsodium-go/cryptosign"
	"github.com/go-errors/errors"
	"github.com/jamesruan/sodium"
	"regexp"
	"strconv"
	"strings"
)

const Ed25519Curve = "Velas seed"
const HardenedOffset = 0x80000000
const defaultDerivePath = "m/0'"

// HD Keypair contains private and public key
type HD struct {
	publicKey  []byte
	privateKey []byte
}

// Generate random keypair, not recommended for using
func GenerateHD() (*HD, error) {
	privateKey, publicKey, err := cryptosign.CryptoSignKeyPair()
	if err != 0 {
		return nil, errors.Errorf("can't generate keys, libsodium error code %d", err)
	}
	return &HD{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

// Generate keypair by private key encrypted in hex
func HDFromPrivateKeyHex(sk string) (*HD, error) {
	bytesSK, err := hex.DecodeString(sk)
	if err != nil {
		return nil, errors.New(err)
	}
	hd := HDFromPrivateKey(bytesSK)
	return &hd, nil
}

// Generate keypair by private key(byte array)
func HDFromPrivateKey(sk []byte) HD {
	ssk := sodium.SignSecretKey{Bytes: sk}
	spk := ssk.PublicKey()

	return HD{
		publicKey:  spk.Bytes,
		privateKey: ssk.Bytes,
	}
}

// Generate keypair by seed and derive index.
func HDFromSeed(seed []byte, pathArg *string) (*HD, error) {
	var path string
	if pathArg == nil {
		path = defaultDerivePath
	} else {
		path = *pathArg
	}
	key, _, err := derivePath(path, seed)
	if err != nil {
		return nil, err
	}
	return hdFromSodiumSeed(key)
}

// generate seed32
func derivePath(path string, seed []byte) ([]byte, []byte, error) {
	if !isValidPath(path) {
		return nil, nil, errors.Errorf("Invalid derivation path %s", path)
	}
	key, chainCode, err := getMasterKeyFromSeed(seed)
	if err != nil {
		return nil, nil, err
	}

	segments, err := pathToSegments(path)
	if err != nil {
		return nil, nil, err
	}

	for _, segment := range segments {
		key, chainCode, err = ckdPriv(key, chainCode, segment+HardenedOffset)
		if err != nil {
			return nil, nil, err
		}
	}
	return key, chainCode, nil
}

func ckdPriv(parentKey []byte, parentChainCode []byte, index uint) ([]byte, []byte, error) {
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, uint32(index))
	data := append(append(make([]byte, 1), parentKey...), indexBytes...)
	h := hmac.New(sha512.New, parentChainCode)
	if _, err := h.Write(data); err != nil {
		return nil, nil, err
	}
	i := h.Sum(nil)
	il := i[:32]
	ir := i[32:]
	return il, ir, nil
}

// parse path string to segments
func pathToSegments(path string) ([]uint, error) {
	arr := pathToStringArray(path)
	segments := make([]uint, 0)
	for _, el := range arr {
		num, err := strconv.Atoi(el)
		if err != nil {
			return nil, errors.New(err)
		}
		segments = append(segments, uint(num))
	}
	return segments, nil
}

func pathToStringArray(path string) []string {
	splitted := strings.Split(path, "/")
	splitted = splitted[1:]
	arr := make([]string, 0)
	// replace derive
	for _, el := range splitted {
		arr = append(arr, strings.Replace(el, "'", "", -1))
	}
	return arr
}

func getMasterKeyFromSeed(seed []byte) ([]byte, []byte, error) {
	h := hmac.New(sha512.New, []byte(Ed25519Curve))
	if _, err := h.Write(seed); err != nil {
		return nil, nil, errors.New(err)
	}
	i := h.Sum(nil)
	il := i[:32]
	ir := i[32:]
	return il, ir, nil
}

// validate path
func isValidPath(path string) bool {
	var re = regexp.MustCompile(`^m(/[0-9]+')+$`)
	return re.MatchString(path)
}

// hd from 32bytes seed
func hdFromSodiumSeed(seed []byte) (*HD, error) {
	sk, pk, errCode := cryptosign.CryptoSignSeedKeyPair(seed)
	if errCode != 0 {
		return nil, errors.Errorf("can't generate keys, libsodium error code %d", errCode)
	}
	return &HD{
		publicKey:  pk,
		privateKey: sk,
	}, nil
}

// Create wallet from hd keypair
func (hd *HD) ToWallet() (*Wallet, error) {
	return CreateWallet(hd.publicKey)
}

// Get private key, encrypted in hex(wif)
func (hd *HD) PrivateKey() string {
	return hex.EncodeToString(hd.privateKey)
}

// Get public key, encrypted in hex
func (hd *HD) PublicKey() string {
	return hex.EncodeToString(hd.publicKey)
}
