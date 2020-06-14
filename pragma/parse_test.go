package pragma

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePragma(t *testing.T) {
	tests := []struct {
		pair      string
		wantKey   string
		wantValue string
		wantOk    bool
	}{
		{
			pair:      "key=value",
			wantKey:   "key",
			wantValue: "value",
			wantOk:    true,
		},
		{
			pair:      "key=value===!",
			wantKey:   "key",
			wantValue: "value===!",
			wantOk:    true,
		},
		{
			pair:      "key",
			wantKey:   "key",
			wantValue: "",
			wantOk:    true,
		},
		{
			pair:      "",
			wantKey:   "",
			wantValue: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s,%t", tt.pair, tt.wantOk), func(t *testing.T) {
			key, value, ok := parsePragma(tt.pair)
			assert.Equal(t, tt.wantKey, key)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}
