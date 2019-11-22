package crypto

import (
	"encoding/hex"
	"github.com/tyler-smith/go-bip39"
	"reflect"
	"testing"
)

func TestHDFromSeed(t *testing.T) {
	const pk = "4766a5dc09364e7f840d291a8883a75dfab09a5c2db5332ffbed56a7ffed5c97310a5225d0c32fbeace30447e12bf2f3d7c629bf563c7aa037fdd803c05371db"
	pkBytes, _ := hex.DecodeString(pk)
	hd := HDFromPrivateKey(pkBytes)
	type args struct {
		mnemonics string
		pathArg   *string
	}
	tests := []struct {
		name    string
		args    args
		want    *HD
		wantErr bool
	}{
		{
			name: "simple seed with default index",
			args: args{
				mnemonics: "grain catch elder liquid ginger daring sure brush sudden whisper garden model",
				pathArg:   nil,
			},
			want:    &hd,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := bip39.NewSeed(tt.args.mnemonics, "")
			got, err := HDFromSeed(seed, tt.args.pathArg)
			if (err != nil) != tt.wantErr {
				t.Errorf("HDFromSeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HDFromSeed() got = %v, want %v", got, tt.want)
			}
		})
	}
}
