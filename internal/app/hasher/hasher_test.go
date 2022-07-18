package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasher_Validate(t *testing.T) {
	h := NewHasher([]byte("test key"))
	data := "test-data"

	encoded, err := h.Sign(data)
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	decoded, err := h.Validate(encoded, int64(len([]byte(data))))

	assert.NoError(t, err)

	assert.Equal(t, data, decoded)
}
