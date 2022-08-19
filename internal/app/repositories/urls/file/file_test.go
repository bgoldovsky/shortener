package file

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repositories/urls/errors"
)

const (
	filePath      = "store.dat"
	defaultUserID = "user123"
)

func TestFileRepo_Add(t *testing.T) {
	ctx := context.Background()
	url := "avito.ru"

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", url, defaultUserID)

	assert.NoError(t, err)
}

func TestFileRepo_Add_Conflict(t *testing.T) {
	ctx := context.Background()
	url := "avito.ru"

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", url, defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "qwerty", url, defaultUserID)
	require.Error(t, err)

	assert.IsType(t, &internalErrors.NotUniqueURLErr{}, err)
}

func TestFileRepo_Get(t *testing.T) {
	ctx := context.Background()
	url := "avito.ru"

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", url, defaultUserID)
	require.NoError(t, err)

	act, err := repo.Get(ctx, "qwerty")

	assert.NoError(t, err)
	assert.Equal(t, "avito.ru", act)
}

func TestFileRepo_Empty(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	act, err := repo.Get(ctx, "qwerty")

	assert.Error(t, err, "url not found")
	assert.Equal(t, "", act)
}

func TestFileRepo_Get_RestoreData(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	repo, err = NewRepository(filePath)
	require.NoError(t, err)

	act, err := repo.Get(ctx, "ytrewq")
	require.NoError(t, err)

	assert.Equal(t, "yandex.ru", act)
}

func TestFileRepo_GetList_Success(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	repo, err = NewRepository(filePath)
	require.NoError(t, err)

	act, err := repo.GetList(ctx, defaultUserID)
	require.NoError(t, err)

	assert.Len(t, act, 2)
}

func TestFileRepository_Delete(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	urlIDs := []string{"qwerty", "ytrewq"}

	err = repo.Add(ctx, urlIDs[0], "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, urlIDs[1], "yandex.ru", defaultUserID)
	require.NoError(t, err)

	repo, err = NewRepository(filePath)
	require.NoError(t, err)

	act, err := repo.GetList(ctx, defaultUserID)
	require.NoError(t, err)
	require.Len(t, act, 2)

	err = repo.Delete(ctx, []models.UserCollection{{UserID: defaultUserID, URLIDs: urlIDs}})
	require.NoError(t, err)

	act, err = repo.GetList(ctx, defaultUserID)
	assert.NoError(t, err)
	assert.Empty(t, act)
}

func TestFileRepo_GetList_NotFound(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Add(ctx, "qwerty", "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	repo, err = NewRepository(filePath)
	require.NoError(t, err)

	act, err := repo.GetList(ctx, "fake")
	require.NoError(t, err)

	assert.Len(t, act, 0)
}

func TestFileRepo_Ping(t *testing.T) {
	ctx := context.Background()

	repo, err := NewRepository(filePath)
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(filePath)
	}()

	err = repo.Ping(ctx)
	assert.NoError(t, err)
}
