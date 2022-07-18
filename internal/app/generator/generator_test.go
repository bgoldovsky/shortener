package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Letters(t *testing.T) {
	tests := []int64{0, 1, 2, 20}
	gen := NewGenerator()

	for _, tt := range tests {
		act, err := gen.RandomString(tt)
		require.NoError(t, err)

		assert.Len(t, act, int(tt))
	}
}
