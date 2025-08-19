package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/kliuchnikovv/packulator/internal/store"
	mock_store "github.com/kliuchnikovv/packulator/internal/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewPackService(t *testing.T) {
	mockStore := mock_store.NewMockStore(gomock.NewController(t))
	service := NewPackService(mockStore)

	assert.NotNil(t, service)
	assert.Implements(t, (*PackService)(nil), service)
}

func TestPackService_CreatePacks(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		packs := []int64{250, 500, 1000}

		// Mock SavePacks to return success
		mockStore.EXPECT().SavePacks(gomock.Any(), gomock.Any()).
			Return(nil)

		versionHash, err := service.CreatePacks(ctx, packs...)

		require.NoError(t, err)
		assert.NotEmpty(t, versionHash)
		assert.Len(t, versionHash, 16) // Hash is truncated to 16 characters
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		packs := []int64{250, 500}
		expectedError := errors.New("database error")

		// Mock SavePacks to return error
		mockStore.EXPECT().SavePacks(gomock.Any(), gomock.Any()).
			Return(expectedError)

		versionHash, err := service.CreatePacks(ctx, packs...)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, versionHash)
	})

	t.Run("empty packs", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		// Mock SavePacks with empty slice
		mockStore.EXPECT().SavePacks(gomock.Any(), gomock.Any()).
			Return(nil)

		versionHash, err := service.CreatePacks(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, versionHash)
	})

	t.Run("single pack", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		// Mock SavePacks
		mockStore.EXPECT().SavePacks(gomock.Any(), gomock.Any()).
			Return(nil)

		versionHash, err := service.CreatePacks(ctx, 1000)

		require.NoError(t, err)
		assert.NotEmpty(t, versionHash)
	})
}

func TestPackService_GetPackByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedPack := &model.Pack{
			ID:          "pack-1",
			VersionHash: "abc123",
			TotalAmount: 250,
			PackItems: []model.PackItem{
				{ID: "item-1", PackID: "pack-1", Size: 250},
			},
		}

		mockStore.EXPECT().GetPackByID(gomock.Any(), "pack-1").Return(expectedPack, nil)

		result, err := service.GetPackByID(ctx, "pack-1")

		require.NoError(t, err)
		assert.Equal(t, expectedPack, result)
	})

	t.Run("pack not found", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		mockStore.EXPECT().GetPackByID(gomock.Any(), "nonexistent").Return(nil, store.ErrNotFound)

		result, err := service.GetPackByID(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, store.ErrNotFound, err)
		assert.Nil(t, result)
	})
}

func TestPackService_ListPacks(t *testing.T) {
	t.Run("successful listing", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedPacks := []model.Pack{
			{
				ID:          "pack-1",
				VersionHash: "abc123",
				TotalAmount: 250,
			},
			{
				ID:          "pack-2",
				VersionHash: "def456",
				TotalAmount: 500,
			},
		}

		mockStore.EXPECT().ListPacks(gomock.Any()).Return(expectedPacks, nil)

		result, err := service.ListPacks(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedPacks, result)
	})

	t.Run("empty list", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		mockStore.EXPECT().ListPacks(gomock.Any()).Return(nil, nil)

		result, err := service.ListPacks(ctx)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedError := errors.New("database error")

		mockStore.EXPECT().ListPacks(gomock.Any()).Return(nil, expectedError)

		result, err := service.ListPacks(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)
	})
}

func TestPackService_DeletePack(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		mockStore.EXPECT().DeletePack(gomock.Any(), "pack-1").Return(nil)

		err := service.DeletePack(ctx, "pack-1")

		require.NoError(t, err)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedError := errors.New("database error")

		mockStore.EXPECT().DeletePack(gomock.Any(), "nonexistent").Return(expectedError)

		err := service.DeletePack(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestGenerateVersionHash(t *testing.T) {
	t.Run("consistent hash for same input", func(t *testing.T) {
		packs := []int64{250, 500, 1000}

		hash1 := generateVersionHash(packs)
		hash2 := generateVersionHash(packs)

		assert.Equal(t, hash1, hash2)
		assert.Len(t, hash1, 16) // Truncated to 16 characters
	})

	t.Run("different hash for different input", func(t *testing.T) {
		packs1 := []int64{250, 500, 1000}
		packs2 := []int64{250, 500, 2000}

		hash1 := generateVersionHash(packs1)
		hash2 := generateVersionHash(packs2)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("same hash for reordered input", func(t *testing.T) {
		packs1 := []int64{250, 500, 1000}
		packs2 := []int64{1000, 250, 500}

		hash1 := generateVersionHash(packs1)
		hash2 := generateVersionHash(packs2)

		assert.Equal(t, hash1, hash2, "hash should be same regardless of input order")
	})

	t.Run("empty input", func(t *testing.T) {
		packs := []int64{}

		hash := generateVersionHash(packs)

		assert.NotEmpty(t, hash)
		assert.Len(t, hash, 16)
	})

	t.Run("single pack", func(t *testing.T) {
		packs := []int64{1000}

		hash := generateVersionHash(packs)

		assert.NotEmpty(t, hash)
		assert.Len(t, hash, 16)
	})

	t.Run("duplicate packs", func(t *testing.T) {
		packs1 := []int64{250, 250, 500}
		packs2 := []int64{500, 250, 250}

		hash1 := generateVersionHash(packs1)
		hash2 := generateVersionHash(packs2)

		assert.Equal(t, hash1, hash2, "hash should handle duplicates correctly")
	})
}

// Benchmark tests
func BenchmarkGenerateVersionHash(b *testing.B) {
	packs := []int64{250, 500, 1000, 2000, 5000}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateVersionHash(packs)
	}
}
