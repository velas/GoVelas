package httpClient

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/velas/GoVelas//httpClient/response"
	"gopkg.in/resty.v1"
)

// Main structure for requesting to node
type Client struct {
	baseAddress string
	Wallet      *Wallet
	Tx          *Tx
	Block       *Block
	bk          *baseClient
}

// Create node client
func NewClient(baseAddress string) *Client {
	bk := newBaseClient(baseAddress)
	return &Client{
		baseAddress: baseAddress,
		bk:          newBaseClient(baseAddress),
		Wallet:      newWalletClient(bk),
		Tx:          newTxClient(bk),
		Block:       newBlockClient(bk),
	}
}

// Method for get status of node
func (cl *Client) NodeInfo() (*response.Node, error) {
	resp, err := resty.
		R().
		Get(cl.baseAddress + "/api/v1/info")
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := cl.bk.ReadResponse(resp)
	result := response.Node{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New(err)
	}
	return &result, nil
}
