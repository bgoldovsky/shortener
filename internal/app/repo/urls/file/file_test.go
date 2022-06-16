package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const filePath = "store.dat"

func TestFileRepo_Add(t *testing.T) {
	repo, err := NewRepo("store.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("qwerty", "avito.ru")
	require.NoError(t, err)
}

func TestFileRepo_Get(t *testing.T) {
	repo, err := NewRepo("store.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("qwerty", "avito.ru")
	require.NoError(t, err)

	act, err := repo.Get("qwerty")
	require.NoError(t, err)

	assert.Equal(t, "avito.ru", act)
}

func TestFileRepo_Empty(t *testing.T) {
	repo, err := NewRepo("store.dat")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	act, err := repo.Get("qwerty")

	assert.Error(t, err, "url not found")
	assert.Equal(t, "", act)
}

func TestFileRepo_Get_RestoreData(t *testing.T) {
	repo, err := NewRepo(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add("qwerty", "avito.ru")
	require.NoError(t, err)

	err = repo.Add("ytrewq", "yandex.ru")
	require.NoError(t, err)

	repo, err = NewRepo("store.dat")
	require.NoError(t, err)

	act, err := repo.Get("ytrewq")
	require.NoError(t, err)

	assert.Equal(t, "yandex.ru", act)
}
