package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repositories/urls/errors"
)

const (
	defaultUserID = "user123"
)

func TestInmemoryRepo_Add(t *testing.T) {
	ctx := context.Background()
	url := "avito.ru"

	repo := NewRepository()

	err := repo.Add(ctx, "qwerty", url, defaultUserID)

	assert.NoError(t, err)
}

func TestInmemoryRepo_Add_Conflict(t *testing.T) {
	ctx := context.Background()
	url := "avito.ru"

	repo := NewRepository()

	err := repo.Add(ctx, "qwerty", url, defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "qwerty", url, defaultUserID)
	require.Error(t, err)

	assert.IsType(t, &internalErrors.NotUniqueURLErr{}, err)
}

func TestInmemoryRepo_Get(t *testing.T) {
	ctx := context.Background()
	url := "avito.ru"

	repo := NewRepository()

	err := repo.Add(ctx, "qwerty", url, defaultUserID)
	require.NoError(t, err)

	act, err := repo.Get(ctx, "qwerty")

	assert.NoError(t, err)
	assert.Equal(t, "avito.ru", act)
}

func TestInmemoryRepo_Empty(t *testing.T) {
	ctx := context.Background()

	repo := NewRepository()
	act, err := repo.Get(ctx, "qwerty")

	assert.Error(t, err, "url not found")
	assert.Equal(t, "", act)
}

func TestInmemoryRepo_GetList_Success(t *testing.T) {
	ctx := context.Background()

	repo := NewRepository()

	err := repo.Add(ctx, "qwerty", "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	act, err := repo.GetList(ctx, defaultUserID)
	require.NoError(t, err)

	assert.Len(t, act, 2)
}

func TestInmemoryRepository_Delete(t *testing.T) {
	ctx := context.Background()

	repo := NewRepository()

	urlIDs := []string{"qwerty", "ytrewq"}

	err := repo.Add(ctx, urlIDs[0], "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, urlIDs[1], "yandex.ru", defaultUserID)
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

func TestInmemoryRepo_GetList_NotFound(t *testing.T) {
	ctx := context.Background()

	repo := NewRepository()

	err := repo.Add(ctx, "qwerty", "avito.ru", defaultUserID)
	require.NoError(t, err)

	err = repo.Add(ctx, "ytrewq", "yandex.ru", defaultUserID)
	require.NoError(t, err)

	act, err := repo.GetList(ctx, "fake")
	require.NoError(t, err)

	assert.Len(t, act, 0)
}

func TestInmemoryRepo_Ping(t *testing.T) {
	ctx := context.Background()

	repo := NewRepository()

	err := repo.Ping(ctx)
	assert.NoError(t, err)
}
