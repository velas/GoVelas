package httpClient

import (
	"reflect"
	"testing"
)

func TestBlock_GetByHash(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *BlockResponse
		wantErr bool
	}{
		{
			name:   "Normal test",
			fields: fields{baseAddress: "https://testnet.velas.com"},
			args: args{
				hash: "508e68a53e8e73e8d7b477e0c3c4dbed83a56899c8dcecabe12f79c42cbd1c87",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Block.GetByHash(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("nil response")
			}
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByHash() got = %v, want %v", got, tt.want)
			}
			*/
			if len(got.Transactions) > 0 && reflect.DeepEqual(got.Transactions[0].Hash, [32]byte{}) {
				t.Error("empty hash")
			}
			if got.Header.TxnCount != uint32(len(got.Transactions)) {
				t.Error("incorrect transactions count")
			}
			t.Logf("%+v", *got)
			t.Logf("%+v", got.Header)
			for _, ad := range got.Advices {
				t.Logf("%+v", ad.PublicKey)
			}
			for _, tx := range got.Transactions {
				t.Logf("%+v", tx)
			}
		})
	}
}
