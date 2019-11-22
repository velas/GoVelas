package httpClient

import (
	"testing"
)

func TestClient_NodeInfo(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Normal test",
			fields:  fields{baseAddress: Url},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.NodeInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("nil response")
			}
			t.Logf("%+v", got)
			t.Logf("P2PInfo %+v", got.P2PInfo)
			t.Logf("Blockchain %+v", got.Blockchain)
			t.Logf("Progress %+v", got.Progress)
		})
	}
}
