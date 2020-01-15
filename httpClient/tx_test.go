package httpClient

import (
	"encoding/hex"
	"github.com/velas/GoVelas/crypto"
	"github.com/velas/GoVelas/crypto/helpers"
	"reflect"
	"testing"
)

func TestTx_GetListByAddress(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		privateKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Normal test",
			fields:  fields{baseAddress: Url},
			args:    args{privateKey: Pk},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			hd, _ := crypto.HDFromPrivateKeyHex(tt.args.privateKey)
			wallet, _ := hd.ToWallet()
			got, err := client.Tx.GetHashListByAddress(wallet.Base58Address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GetListByAddress() got empty array")
			}
			t.Log(got)
		})
	}
}

// Test doesn't work correct
func TestTx_GetByHashList(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		hashes []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TxResponse
		wantErr bool
	}{
		{
			name:   "Normal test",
			fields: fields{baseAddress: Url},
			args: args{hashes: []string{
				"de7efc5dff6860bdb78758a4851c48ab284c039b85320a7dd334648cb787e317",
				"ca4161d7743a93d4a1c5c4ba8462435dc1a23219941fea5cebc26ff93d2bcec6",
			}},
			want: []TxResponse{
				{
					Size:               39,
					Block:              "519cb2a564e0e682540663af09ef56f9068aeb5c51f8f852a3011ff860a479f9",
					Confirmed:          7745,
					ConfirmedTimestamp: 1565085672,
					Total:              0,
					Tx:                 nil,
				}, {
					Size:               38,
					Block:              "498ae932b8066aefa14596547111fa957e9a42dd39cbac993564df2e1102f280",
					Confirmed:          85229,
					ConfirmedTimestamp: 1563535015,
					Total:              0,
					Tx:                 nil,
				},
			},
			wantErr: false,
		},
		{
			name:    "check tx confirmed",
			args:    args{hashes: []string{"53ab5f62deac40f68e18c0600e775c95ccde6b6fa7e9bc552d6109431571896f"}},
			fields:  fields{baseAddress: Url},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Tx.GetByHashList(tt.args.hashes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetListByAddress() got = %v, want %v", got, tt.want)
			}*/
			t.Logf("%+v", got)
		})
	}
}

func TestTx_Validate(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		privateKey string
		toAddress  string
		amount     uint64
		commission uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TxResponse
		wantErr bool
	}{
		{
			name:   "Normal test",
			fields: fields{baseAddress: Url},
			args: args{
				privateKey: Pk,
				toAddress:  "VLa1hi77ZXD2BSWDD9wQe8vAhejXyS7vBM4",
				amount:     1000,
				commission: 1000000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			nodeIDBytes, _ := hex.DecodeString("f77da461bf9c71df996b4b33c54fed6877179bc33656d3aa62e44782bf4eb1f7")
			nodeID := helpers.ToHash(nodeIDBytes)
			hd, _ := crypto.HDFromPrivateKeyHex(tt.args.privateKey)
			wallet, _ := hd.ToWallet()
			unspents, _ := client.Wallet.GetUnspent(wallet.Base58Address)
			tx, err := crypto.NewTransaction(unspents, tt.args.amount, *hd, wallet.Base58Address, tt.args.toAddress, tt.args.commission, nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tx == nil {
				t.Errorf("Publish() error = %v, wantErr %v", " tx is nil", tt.wantErr)
				return
			}
			err = client.Tx.Validate(*tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(hex.EncodeToString(tx.Hash[:]))
		})
	}
}

func TestTx_Publish(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		privateKey   string
		toAddress    string
		amount       uint64
		commission   uint64
		nodeIDString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TxResponse
		wantErr bool
	}{
		{
			name:   "Normal test",
			fields: fields{baseAddress: TestNetUrl},
			args: args{
				privateKey:   Pk,
				toAddress:    "VLa1hi77ZXD2BSWDD9wQe8vAhejXyS7vBM4",
				amount:       14,
				commission:   1000000,
				nodeIDString: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			hd, _ := crypto.HDFromPrivateKeyHex(tt.args.privateKey)
			wallet, _ := hd.ToWallet()
			unspents, _ := client.Wallet.GetUnspent(wallet.Base58Address)
			nodeIDBytes, _ := hex.DecodeString(tt.args.nodeIDString)
			var nodeID crypto.NodeID
			if len(nodeIDBytes) == 32 {
				nodeID = helpers.ToHash(nodeIDBytes)
			} else {
				nodeID = [32]byte{}
			}
			tx, err := crypto.NewTransaction([]crypto.TransactionInputOutpoint{unspents[0]}, tt.args.amount, *hd, wallet.Base58Address, tt.args.toAddress, tt.args.commission, nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tx == nil {
				t.Errorf("Publish() error = %v, wantErr %v", " tx is nil", tt.wantErr)
				return
			}
			err = client.Tx.Publish(*tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(hex.EncodeToString(tx.Hash[:]))
		})
	}
}

func TestTx_MakeStake(t *testing.T) {
	hd, _ := crypto.HDFromPrivateKeyHex(Pk2)
	wallet, _ := hd.ToWallet()
	returnAddress := wallet.Base58Address
	staker, _ := hex.DecodeString("ccf62994d19a2a08fc8952ad90c2f20830bfad5ba2779bc03c65aa51a7c4c1fc")
	fakeStaker, _ := hex.DecodeString("f77da461bf9c71df996b4b33c54fed6877179bc33656d3aa62e44782bf4eb1f7")
	secondStaker, _ := hex.DecodeString("7779fa79d07aaefb5412b2b51abed9be8cafc4cf2961192aee3d573166e8a168")
	secondStakerPk := "2a53e9463e66ce9e0589ab7a93c2626550d9e602cdf0062126bc54ac1dd774347779fa79d07aaefb5412b2b51abed9be8cafc4cf2961192aee3d573166e8a168"
	ssHD, _ := crypto.HDFromPrivateKeyHex(secondStakerPk)
	ssWallet, _ := ssHD.ToWallet()
	type fields struct {
		baseAddress string
	}
	type args struct {
		privateKey string
		commission uint64
		receivers  []crypto.Receiver
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TxResponse
		wantErr bool
	}{
		{
			name:   "Staking test",
			fields: fields{baseAddress: LocalUrl},
			args: args{
				privateKey: Pk2,
				commission: 1000000,
				receivers: []crypto.Receiver{
					{
						Wallet: returnAddress,
						Amount: 1000000000000,
						NodeID: helpers.ToHash(staker),
					},
					{
						Wallet: returnAddress,
						Amount: 2000000000000,
						NodeID: helpers.ToHash(fakeStaker),
					},
					{
						Wallet: ssWallet.Base58Address,
						Amount: 1900000000000,
						NodeID: helpers.ToHash(secondStaker),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			hd, _ := crypto.HDFromPrivateKeyHex(tt.args.privateKey)
			wallet, _ := hd.ToWallet()
			unspents, _ := client.Wallet.GetUnspentForStaking(wallet.Base58Address)
			tx, err := crypto.NewTransactionManyRecievers(
				unspents, *hd, wallet.Base58Address, tt.args.receivers, tt.args.commission)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tx == nil {
				t.Errorf("Publish() error = %v, wantErr %v", " tx is nil", tt.wantErr)
				return
			}
			err = client.Tx.Publish(*tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(hex.EncodeToString(tx.Hash[:]))
		})
	}
}

func TestTx_GetHashListByHeight(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		height int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Correct",
			fields:  fields{baseAddress: Url},
			args:    args{height: 209126},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Tx.GetHashListByHeight(tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHashListByHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHashListByHeight() got = %v, want %v", got, tt.want)
			}
		})
	}
}
