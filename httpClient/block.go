package httpClient

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/velas/GoVelas/crypto"
	"gopkg.in/resty.v1"
)

// Block node client
type Block struct {
	bk *baseClient
}

// create block node client
func newBlockClient(bk *baseClient) *Block {
	return &Block{
		bk: bk,
	}
}

// Block object response from node request
type BlockResponse struct {
	Header       *Header      `json:"header"`  // Block header
	Transactions []*crypto.Tx `json:"txns"`    // Block transactions, in format of "tx" command
	Advices      []*Advice    `json:"advices"` // list of advicers, who can sign the blocks
}

type Advice struct {
	PublicKey string `json:"public_key"`
}

// Header object response from node request
type Header struct {
	Type        uint32 `json:"type"`         // Block type
	Hash        string `json:"hash"`         // Hash
	Height      uint32 `json:"height"`       // Height
	Size        uint64 `json:"size"`         // Size
	Version     uint32 `json:"version"`      // Block version information (note, this is signed)
	PrevBlock   string `json:"prev_block"`   // The hash value of the previous block this particular block references
	MerkleRoot  string `json:"merkle_root"`  // The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp   uint32 `json:"timestamp"`    // A timestamp recording when this block was created (Will overflow in 2106[2])
	Bits        uint32 `json:"bits"`         // Not used
	Nonce       uint32 `json:"nonce"`        // The nonce used to generate this block… to allow variations of the header and compute different hashes
	Seed        string `json:"seed"`         // The random seed, not used
	TxnCount    uint32 `json:"txn_count"`    // Transaction count
	AdviceCount uint32 `json:"advice_count"` // Advise list count
	Script      string `json:"script"`       // The node's (block owner) signature
}

// Method for get block object
func (blk *Block) GetByHash(hash string) (*BlockResponse, error) {
	resp, err := resty.
		R().
		Get(blk.bk.baseAddress + "/api/v1/blocks/" + hash)
	if err != nil {
		return nil, err
	}
	body, err := blk.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	blockResponse := BlockResponse{}
	if err := json.Unmarshal(body, &blockResponse); err != nil {
		return nil, errors.New(err)
	}
	return &blockResponse, nil
}
