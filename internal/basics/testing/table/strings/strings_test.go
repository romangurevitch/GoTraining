package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrings_AddPrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid prefix",
			prefix:  "pre_",
			input:   "hello",
			want:    "pre_hello",
			wantErr: false,
		},
		{
			name:    "empty input returns error",
			prefix:  "pre_",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty prefix",
			prefix:  "",
			input:   "hello",
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "multi-word prefix",
			prefix:  "hello world ",
			input:   "foo",
			want:    "hello world foo",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Strings{prefix: tt.prefix}
			got, err := s.AddPrefix(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
