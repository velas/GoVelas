package crypto

import (
	"reflect"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	type args struct {
		privateKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Simple",
			args:    args{privateKey: "d4e3d9ee7c9f57dc35db5dd2360f6ee8b5085a9a14878e234b38f154e18b449bd352b77fd8c136ecb42a0909f9496bbaf3229ba243190488bd1fbf9fa62a7ec5"},
			want:    "VLcN1dBy1VPc9bijr8rzeGbC78MCQ8DjwvS",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd, err := HDFromPrivateKeyHex(tt.args.privateKey)
			if err != nil {
				t.Error(err)
			}
			got, err := hd.ToWallet()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Base58Address, tt.want) {
				t.Errorf("CreateWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
