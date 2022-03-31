package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_increment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tag        string
		want       string
		wantPrefix string
	}{
		{
			name:       "Should return expected",
			tag:        "v0.0.1",
			want:       "v0.0.2",
			wantPrefix: "",
		},
		{
			name:       "Should return expected, v1.1.110",
			tag:        "v1.1.110",
			want:       "v1.1.111",
			wantPrefix: "",
		},
		{
			name:       "Should return expected, hoge/v0.0.2",
			tag:        "hoge/microservice/v0.0.2",
			want:       "hoge/microservice/v0.0.3",
			wantPrefix: "hoge/microservice/",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag, p := increment(tt.tag)
			assert.Equal(t, tt.want, tag)
			assert.Equal(t, tt.wantPrefix, p)
		})
	}
}

func Test_writeReleaseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filePath string
		tag      string
		wantTag  string
	}{
		{
			name: "increment release file",
		},
		{
			name: "replacement release file",
			tag:  "v1.1.1",
		},
		{
			name:     "replacement release file, has prefix",
			filePath: "./testdata/RELEASE",
			tag:      "hoge/fuga/v1.1.1",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			writeReleaseFile(tt.filePath, tt.tag)
		})
	}
}
