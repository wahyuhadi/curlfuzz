package curlfuzz

import "testing"

func TestCurlToNuc(t *testing.T) {
	type args struct {
		c   string
		key string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testing 1",
			args: args{
				c:   "curl -d 'key=aaa' http://127.0.0.1",
				key: "wa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CurlToNuc(tt.args.c, tt.args.key)
		})
	}
}
