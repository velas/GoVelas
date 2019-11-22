package httpClient

import (
	"encoding/hex"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/velas/GoVelas//crypto"
	"gopkg.in/resty.v1"
	"strconv"
)

// Transaction response from node
type TxResponse struct {
	Size               uint32 `json:"size"`
	Block              string `json:"block"` // block hash
	Confirmed          uint32 `json:"confirmed"`
	ConfirmedTimestamp uint32 `json:"confirmed_timestamp"`
	Total              int    `json:"total,omitempty"`
	*crypto.Tx
}

// Custom unmarshaller of transaction response
func (txr *TxResponse) UnmarshalJSON(data []byte) error {
	type Alias TxResponse
	aux := &struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Alias: (*Alias)(txr),
	}
	if aux.Tx == nil {
		aux.Alias.Tx = &crypto.Tx{}
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	pHash, err := hex.DecodeString(aux.Hash)
	if err != nil {
		return err
	}

	var hash [32]byte
	if len(pHash) == 32 {
		copy(hash[:], pHash[:32])
	}
	txr.Hash = hash
	return nil
}

// Transaction client
type Tx struct {
	bk *baseClient
}

// create transaction client
func newTxClient(bk *baseClient) *Tx {
	return &Tx{
		bk: bk,
	}
}

// Get an array of transaction hashes by wallet address
func (tx *Tx) GetHashListByAddress(address string) ([]string, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/v1/wallet/txs/" + address)
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	response := make([]string, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response, nil
}

// Get an array of transaction hashes by range blocks from the given height to the highest block
func (tx *Tx) GetHashListByHeight(height int) ([]string, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/v1/txs/height/" + strconv.Itoa(height))
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	response := make([]string, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response, nil
}

// Get array of transaction objects by hash list, maximum hashes is 10000(can change later)
func (tx *Tx) GetByHashList(hashes []string) ([]TxResponse, error) {
	arg := struct {
		Hashes []string `json:"hashes"`
	}{Hashes: hashes}
	resp, err := resty.
		R().
		SetBody(arg).
		Post(tx.bk.baseAddress + "/api/v1/txs")
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	response := make([]TxResponse, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response, nil
}

// Response on publish and validate requests
type TxPublishResponse struct {
	Result string `json:"result"`
}

// Method for validate transaction, return error if transaction incorrect. Does not publish transactions on the
// blockchain
func (tx *Tx) Validate(txData crypto.Tx) error {
	resp, err := resty.
		R().
		SetBody(&txData).
		Post(tx.bk.baseAddress + "/api/v1/txs/validate")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return errors.New(err)
	}
	return nil
}

// Publish transaction in blockchain
func (tx *Tx) Publish(txData crypto.Tx) error {
	resp, err := resty.
		R().
		SetBody(&txData).
		Post(tx.bk.baseAddress + "/api/v1/txs/publish")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return errors.New(err)
	}
	return nil
}
