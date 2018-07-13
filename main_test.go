package main

import "testing"

func Test_generateAppURL(t *testing.T) {

	tests := []struct {
		name      string
		publicKey string
		want      string
	}{
		{
			name:      "test_1",
			publicKey: "u2c556dxrxjfjzyh0xgpvb1tvr",
			want:      "https://appetize.io/app/u2c556dxrxjfjzyh0xgpvb1tvr",
		},
		{
			name:      "test_2",
			publicKey: "u2c556dxrxjfjzyh0xgpvb1tvr",
			want:      "https://appetize.io/app/u2c556dxrxjfjzyh0xgpvb1tvr",
		},
		{
			name:      "test_3",
			publicKey: "u2c556dxrxjfjzyh0xgpvb1tvr",
			want:      "https://appetize.io/app/u2c556dxrxjfjzyh0xgpvb1tvr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateAppURL(tt.publicKey); got != tt.want {
				t.Errorf("generateAppURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
