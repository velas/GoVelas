package crypto

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/velas/GoVelas//crypto/helpers"
)

// Transaction input
type TransactionInput struct {
	// The previous output transaction reference, as an OutPoint structure
	PreviousOutput TransactionInputOutpoint `json:"previous_output"`
	// Transaction version as defined by the sender. Intended for "replacement" of transactions when information is
	// updated before inclusion into a block.
	Sequence      uint32 `json:"sequence"`
	Script        []byte `json:"signature_script"`     // Computational Script for confirming transaction authorization
	PublicKey     []byte `json:"public_key,omitempty"` // Public Key for verify signature when to work with wallet
	WalletAddress []byte `json:"wallet_address,omitempty"`
}

// Previous transaction out
type TransactionInputOutpoint struct {
	Hash  [32]byte `json:"hash"`  // The hash of the referenced transaction
	Index uint32   `json:"index"` // The index of the specific output in the transaction. The first output is 0, etc.
	Value uint64   `json:"value"` // Transaction Value
}

// ToBytes convert TransactionInputOutpoint to bytes slice
func (tio TransactionInputOutpoint) ToBytes() []byte {
	slices := [][]byte{
		tio.Hash[:],                      // 32 bytes
		helpers.UInt32ToBytes(tio.Index), // 4 bytes
		helpers.UInt64ToBytes(tio.Value), // 8 bytes
	}

	return helpers.ConcatByteArray(slices)
}

// MarshalJSON custom json convert
func (tio *TransactionInputOutpoint) MarshalJSON() ([]byte, error) {
	type Alias TransactionInputOutpoint
	return json.Marshal(&struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Hash:  hex.EncodeToString(tio.Hash[:]),
		Alias: (*Alias)(tio),
	})
}

// UnmarshalJSON custom json convert
func (tio *TransactionInputOutpoint) UnmarshalJSON(data []byte) error {
	type Alias TransactionInputOutpoint
	aux := &struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Alias: (*Alias)(tio),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	hashStop, err := hex.DecodeString(aux.Hash)
	if err != nil {
		return err
	}

	var hash [32]byte
	copy(hash[:], hashStop[:32])
	tio.Hash = hash
	return nil
}

// Transaction input to byte array for generate hash
func (ti *TransactionInput) forBlkHash() []byte {
	slices := [][]byte{
		ti.PreviousOutput.ToBytes(),        // SCTxInOutopointLen 44
		helpers.UInt32ToBytes(ti.Sequence), // 4 bytes
		ti.PublicKey,                       // 32 bytes
		[]byte(ti.Script),                  // sc.ScriptLength
	}

	return helpers.ConcatByteArray(slices)
}

// Marshal Transaction input to JSON byte array
func (ti *TransactionInput) MarshalJSON() ([]byte, error) {
	type Alias TransactionInput
	return json.Marshal(&struct {
		Script        string `json:"signature_script"`
		WalletAddress string `json:"wallet_address"`
		*Alias
	}{
		Script:        hex.EncodeToString(ti.Script),
		WalletAddress: base58.Encode(ti.WalletAddress),
		Alias:         (*Alias)(ti),
	})
}

// UnmarshalJSON custom json convert
func (ti *TransactionInput) UnmarshalJSON(data []byte) error {
	type Alias TransactionInput
	aux := &struct {
		Script        string `json:"signature_script"`
		WalletAddress string `json:"wallet_address,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(ti),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error

	ti.Script, err = hex.DecodeString(aux.Script)
	if err != nil {
		return err
	}
	ti.WalletAddress = base58.Decode(aux.WalletAddress)

	return nil
}
