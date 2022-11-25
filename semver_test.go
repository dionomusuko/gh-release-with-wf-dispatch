package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSemver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		tag       string
		incrLevel string
		want      string
		wantErr   string
	}{
		{
			name:      "success with patch",
			tag:       "v1.0.0",
			incrLevel: "patch",
			want:      "v1.0.1",
		},
		{
			name:      "success with minor",
			tag:       "v1.0.0",
			incrLevel: "minor",
			want:      "v1.1.0",
		},
		{
			name:      "success with major",
			tag:       "v1.0.0",
			incrLevel: "major",
			want:      "v2.0.0",
		},
		{
			name:      "success with monorepo tag format",
			tag:       "app1/microservice/v1.0.0",
			incrLevel: "patch",
			want:      "app1/microservice/v1.0.1",
		},
		{
			name:      "success with YAML comment",
			tag:       "app1/microservice/v1.0.0 # comment",
			incrLevel: "patch",
			want:      "app1/microservice/v1.0.1",
		},
		{
			name:      "failed to parse semver",
			tag:       "va.b.c",
			incrLevel: "patch",
			want:      "",
			wantErr:   "invalid semantic version",
		},
		{
			name:      "failed with empty level",
			tag:       "v1.0.0",
			incrLevel: "",
			want:      "",
			wantErr:   "next_semver_level is empty",
		},
		{
			name:      "failed with invalid level",
			tag:       "v1.0.0",
			incrLevel: "invalid",
			want:      "",
			wantErr:   "invalid not supported",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			newSemver, err := newSemver(tt.tag, tt.incrLevel)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
			assert.Equal(t, tt.want, newSemver)
		})
	}
}
