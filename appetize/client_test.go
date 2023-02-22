package appetize

import "testing"

func Test_baseURL(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		appPath     string
		publicKey   string
		apiEndPoint string
		want        string
	}{
		{
			name:      "test_1",
			token:     "token_abcdefg",
			appPath:   "./apps/XcodeArchiveTest.app.zip",
			publicKey: "",
			apiEndPoint: "@api.appetize.io/v1/apps",
			want:      "https://token_abcdefg@api.appetize.io/v1/apps",
		},
		{
			name:      "test_2",
			token:     "token_abcdefg",
			appPath:   "./apps/XcodeArchiveTest.app.zip",
			publicKey: "pubkey",
			apiEndPoint: "@api.appetize.io/v1/apps",
			want:      "https://token_abcdefg@api.appetize.io/v1/apps/pubkey",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := baseURL(tt.token, tt.appPath, tt.publicKey,tt.apiEndPoint); got != tt.want {
				t.Errorf("baseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
