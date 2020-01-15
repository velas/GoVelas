package httpClient

import (
	"github.com/velas/GoVelas/crypto"
	"testing"
)

const TestNetUrl = "https://testnet.velas.com"
const TestNetMasterNodeUrl = "https://testnet.velas.com"
const Url = TestNetUrl
const LocalUrl = "http://localhost:8088"

const Pk = "89d5bd2d31889df63cb1c895e4c6f16772e7b06a8c71228bb59d4c9a0c434fc1f6e586d5d051065a580969d15f48f88251ed24b9c77422410bc39a0e7247e53a"
const Pk2 = "caa5802c315c994651e757ab5ae2de1f087ba4588e30cffb3fe7ac022ba4ecc6e6bb0082a92e91f92a5480a1f5d4df435f6752b4b31b3d06c11d126a98bfd978"

func TestGetWalletBalance(t *testing.T) {
	client := NewClient(TestNetMasterNodeUrl)
	hd, err := crypto.HDFromPrivateKeyHex(Pk2)
	if err != nil {
		t.Error(err)
	}
	wallet, err := hd.ToWallet()
	if err != nil {
		t.Error(err)
	}
	balance, err := client.Wallet.GetBalance(wallet.Base58Address)
	if err != nil {
		t.Error(err)
	}
	t.Log(balance)
}

func TestGetWalletUnspents(t *testing.T) {
	client := NewClient(TestNetMasterNodeUrl)
	hd, err := crypto.HDFromPrivateKeyHex(Pk2)
	if err != nil {
		t.Error(err)
	}
	wallet, err := hd.ToWallet()
	if err != nil {
		t.Error(err)
	}
	unspents, err := client.Wallet.GetUnspent(wallet.Base58Address)
	if err != nil {
		t.Error(err)
	}
	t.Log(unspents)
}
