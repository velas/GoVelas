package crypto

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/velas/GoVelas//crypto/helpers"
)

// Transaction output
type TransactionOutput struct {
	Index         uint32 `json:"index"`
	Value         uint64 `json:"value"`     // Transaction Value
	Script        []byte `json:"pk_script"` // Usually contains the public key as a script setting up conditions to claim this output.
	Payload       []byte `json:"payload"`
	WalletAddress []byte `json:"wallet_address,omitempty"`
	NodeID        NodeID `json:"node_id"` // usually contains the public key of node for request to be part of commission
}

// Node ID use only for stacking transaction
type NodeID [32]byte

// IsEmpty check script is empty
func (n NodeID) IsEmpty() bool {
	empty := [32]byte{}
	ch := n[:]
	if bytes.Equal(empty[:], ch) {
		return true
	}
	if len(ch) == 0 {
		return true
	}
	return false
}

// forBlkHash - convert transaction output to byte slice
func (to *TransactionOutput) forBlkHash() []byte {
	slices := [][]byte{
		helpers.UInt32ToBytes(to.Index), // 4 bytes
		helpers.UInt64ToBytes(to.Value), // 8 bytes
		to.Script,                       // Script
	}
	if !to.NodeID.IsEmpty() {
		slices = append(slices, to.NodeID[:])
	}

	return helpers.ConcatByteArray(slices)
}

// Generate message for sign out
func (to *TransactionOutput) msgForSign() []byte {
	slices := [][]byte{
		helpers.UInt32ToBytes(to.Index), // 4 bytes
		helpers.UInt64ToBytes(to.Value), // 8 bytes
		to.Script,                       // sc.ScriptLength
	}

	if !to.NodeID.IsEmpty() {
		slices = append(slices, to.NodeID[:])
	}

	return helpers.ConcatByteArray(slices)
}

// MarshalJSON custom json convert
func (to *TransactionOutput) MarshalJSON() ([]byte, error) {
	type Alias TransactionOutput
	return json.Marshal(&struct {
		Script        string `json:"pk_script"`
		WalletAddress string `json:"wallet_address,omitempty"`
		NodeID        string `json:"node_id"`
		*Alias
	}{
		Script:        hex.EncodeToString(to.Script),
		NodeID:        hex.EncodeToString(to.NodeID[:]),
		WalletAddress: base58.Encode(to.Script),
		Alias:         (*Alias)(to),
	})
}

// UnmarshalJSON custom json convert
func (to *TransactionOutput) UnmarshalJSON(data []byte) error {
	type Alias TransactionOutput
	aux := &struct {
		Script        string `json:"pk_script"`
		NodeID        string `json:"node_id"`
		WalletAddress string `json:"wallet_address,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(to),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error

	to.Script, err = hex.DecodeString(aux.Script)
	if err != nil {
		return err
	}

	nodeIDBuf, err := hex.DecodeString(aux.NodeID)
	if err != nil {
		return err
	}
	to.NodeID, err = helpers.GetHash(nodeIDBuf)
	if err != nil {
		return err
	}
	to.WalletAddress = base58.Decode(aux.WalletAddress)
	return nil
}
