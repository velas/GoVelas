package httpClient

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/velas/GoVelas//crypto"
	"gopkg.in/resty.v1"
)

// Wallet node client
type Wallet struct {
	bk *baseClient
}

// Balance response
type Balance struct {
	Amount uint64 `json:"amount"`
}

// Create wallet node client
func newWalletClient(bk *baseClient) *Wallet {
	return &Wallet{
		bk: bk,
	}
}

// Get balance of wallet by base58 address
func (w *Wallet) GetBalance(address string) (uint64, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/v1/wallet/balance/" + address)
	if err != nil {
		return 0, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return 0, err
	}
	balanceResponse := Balance{}
	if err := json.Unmarshal(body, &balanceResponse); err != nil {
		return 0, errors.New(err)
	}
	return balanceResponse.Amount, nil
}

// Get unspent outs for using in transactions
func (w *Wallet) GetUnspent(address string) ([]crypto.TransactionInputOutpoint, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/v1/wallet/unspent/" + address)
	if err != nil {
		return nil, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	unspents := make([]crypto.TransactionInputOutpoint, 0)
	if err := json.Unmarshal(body, &unspents); err != nil {
		return nil, errors.New(err)
	}
	return unspents, nil
}

func (w *Wallet) GetUnspentForStaking(address string) ([]crypto.TransactionInputOutpoint, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/v1/wallet/unspent_for_staking/" + address)
	if err != nil {
		return nil, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	unspents := make([]crypto.TransactionInputOutpoint, 0)
	if err := json.Unmarshal(body, &unspents); err != nil {
		return nil, errors.New(err)
	}
	return unspents, nil
}
