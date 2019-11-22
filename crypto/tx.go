package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/GoKillers/libsodium-go/cryptosign"
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-errors/errors"
	"github.com/velas/GoVelas/crypto/helpers"
)

// Transaction object
type Tx struct {
	Hash     [32]byte            `json:"hash"`
	Version  uint32              `json:"version"`
	LockTime uint32              `json:"lock_time"`
	Inputs   []TransactionInput  `json:"tx_in"`
	Outputs  []TransactionOutput `json:"tx_out"`
}

// To create a transaction, it is necessary to receive unspent outputs of the sender’s wallet, all outs will be spent on
// commission, to the recipient and the remaining amount will be returned to the sender. Sign method generate hash of
// transaction.
func NewTransaction(
	unspents []TransactionInputOutpoint,
	amount uint64,
	key HD,
	fromAddress string,
	to string,
	commission uint64,
	nodeID NodeID,
) (*Tx, error) {
	totalin := int64(0)

	for _, previousOutput := range unspents {
		totalin += int64(previousOutput.Value)
	}

	index := uint32(0)

	txIns := make([]TransactionInput, 0)

	txOuts := make([]TransactionOutput, 0)

	txOuts = append(txOuts, TransactionOutput{
		Index: index,
		Value: commission,
	})

	index++
	txOuts = append(txOuts, TransactionOutput{
		Index:         index,
		Script:        base58.Decode(to),
		Value:         amount,
		WalletAddress: base58.Decode(to),
		// NodeID:        nodeID,
	})

	change := totalin - int64(amount) - int64(commission)

	if change < 0 {
		return nil, errors.Errorf("Insufficient funds, total amount %d, commission %d, send amount %d", totalin, commission, amount)
	} else if change > 0 {
		// My address
		index++
		txOuts = append(txOuts, TransactionOutput{
			Index:         index,
			Script:        base58.Decode(fromAddress),
			Value:         uint64(change),
			WalletAddress: base58.Decode(fromAddress),
			NodeID:        nodeID,
		})
	}

	tx := Tx{
		Version:  1,
		LockTime: 0,
		Outputs:  txOuts,
	}

	for _, previousOutput := range unspents {
		sigMsg := tx.msgForSign(previousOutput.Hash, previousOutput.Index)
		sig, err := cryptosign.CryptoSignDetached(sigMsg, key.privateKey)
		if err != 0 {
			return nil, errors.Errorf("Error on sign message")
		}
		txIns = append(txIns, TransactionInput{
			PublicKey:      key.publicKey,
			Sequence:       1,
			PreviousOutput: previousOutput,
			Script:         sig,
			WalletAddress:  base58.Decode(fromAddress),
		})
	}
	tx.Inputs = txIns
	txHash := tx.generateHash()
	tx.Hash = txHash
	return &tx, nil
}

type Receiver struct {
	Wallet string
	Amount uint64
	NodeID NodeID
}

func NewTransactionManyRecievers(
	unspents []TransactionInputOutpoint,
	key HD,
	fromAddress string,
	receivers []Receiver,
	commission uint64,
) (*Tx, error) {
	totalin := int64(0)
	totalout := int64(0)

	for _, previousOutput := range unspents {
		totalin += int64(previousOutput.Value)
	}

	for _, receiver := range receivers {
		totalout += int64(receiver.Amount)
	}

	index := uint32(0)

	txIns := make([]TransactionInput, 0)

	txOuts := make([]TransactionOutput, 0)

	txOuts = append(txOuts, TransactionOutput{
		Index: index,
		Value: commission,
	})

	for _, receiver := range receivers {
		index++
		txOuts = append(txOuts, TransactionOutput{
			Index:         index,
			Script:        base58.Decode(receiver.Wallet),
			Value:         receiver.Amount,
			WalletAddress: base58.Decode(receiver.Wallet),
			NodeID:        receiver.NodeID,
		})
	}

	change := totalin - totalout - int64(commission)

	if change < 0 {
		return nil, errors.Errorf("Insufficient funds, total amount %d, commission %d, send amount %d", totalin, commission, totalout)
	} else if change > 0 {
		// My address
		index++
		txOuts = append(txOuts, TransactionOutput{
			Index:         index,
			Script:        base58.Decode(fromAddress),
			Value:         uint64(change),
			WalletAddress: base58.Decode(fromAddress),
		})
	}

	tx := Tx{
		Version:  1,
		LockTime: 0,
		Outputs:  txOuts,
	}

	for _, previousOutput := range unspents {
		sigMsg := tx.msgForSign(previousOutput.Hash, previousOutput.Index)
		sig, err := cryptosign.CryptoSignDetached(sigMsg, key.privateKey)
		if err != 0 {
			return nil, errors.Errorf("Error on sign message")
		}
		txIns = append(txIns, TransactionInput{
			PublicKey:      key.publicKey,
			Sequence:       1,
			PreviousOutput: previousOutput,
			Script:         sig,
			WalletAddress:  base58.Decode(fromAddress),
		})
	}
	tx.Inputs = txIns
	txHash := tx.generateHash()
	tx.Hash = txHash
	return &tx, nil
}

// msgForSign return msg for sign transaction inputs
func (tx *Tx) msgForSign(hash [32]byte, index uint32) []byte {
	txOutSlices := make([][]byte, 0)
	for _, txOut := range tx.Outputs {
		txOutSlices = append(txOutSlices, txOut.msgForSign())
	}
	txOutSlice := helpers.ConcatByteArray(txOutSlices)

	txSlices := [][]byte{
		hash[:],                            // 32
		helpers.UInt32ToBytes(index),       // 4 bytes
		helpers.UInt32ToBytes(tx.Version),  // 4 bytes
		helpers.UInt32ToBytes(tx.LockTime), // 4 bytes
		txOutSlice,
	}

	return helpers.ConcatByteArray(txSlices)
}

// generateHash return generated hash
func (tx *Tx) generateHash() [32]byte {
	txInSlices := make([][]byte, 0)
	for _, txIn := range tx.Inputs {
		txInSlices = append(txInSlices, txIn.forBlkHash())
	}
	txInSlice := helpers.ConcatByteArray(txInSlices)

	txOutSlices := make([][]byte, 0)
	for _, txOut := range tx.Outputs {
		txOutSlices = append(txOutSlices, txOut.forBlkHash())
	}
	txOutSlice := helpers.ConcatByteArray(txOutSlices)

	txSlices := [][]byte{
		helpers.UInt32ToBytes(tx.Version),  // 4 bytes
		helpers.UInt32ToBytes(tx.LockTime), // 4 bytes
		txInSlice,
		txOutSlice,
	}

	msg := helpers.ConcatByteArray(txSlices)
	return DHASH(msg)
}

// Double sha256 hash
func DHASH(data []byte) [32]byte {
	sum := sha256.Sum256(data)
	sum = sha256.Sum256(sum[:])
	return sum
}

// Marshal Tx to json string
func (tx *Tx) MarshalJSON() ([]byte, error) {
	type Alias Tx
	return json.Marshal(&struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Hash:  hex.EncodeToString(tx.Hash[:]),
		Alias: (*Alias)(tx),
	})
}

// UnmarshalJSON custom json convert
func (tx *Tx) UnmarshalJSON(data []byte) error {
	type Alias Tx
	aux := &struct {
		Hash string `json:"hash"`
		*Alias
	}{
		Alias: (*Alias)(tx),
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
	tx.Hash = hash
	return nil
}
