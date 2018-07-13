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
			publicKey: "u2c556dxrxjfjzy",
			want:      "https://appetize.io/app/u2c556dxrxjfjzy",
		},
		{
			name:      "test_2",
			publicKey: "h0xgpvb1tvr",
			want:      "https://appetize.io/app/h0xgpvb1tvr",
		},
		{
			name:      "test_3",
			publicKey: "lkwo92okosss",
			want:      "https://appetize.io/app/lkwo92okosss",
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
