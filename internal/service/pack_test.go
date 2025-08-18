package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/kliuchnikovv/packulator/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStore implements store.Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetPacksInvariantsByHash(ctx context.Context, versionHash string) ([]model.Pack, error) {
	args := m.Called(ctx, versionHash)
	return args.Get(0).([]model.Pack), args.Error(1)
}

func (m *MockStore) SavePack(ctx context.Context, pack *model.Pack) error {
	args := m.Called(ctx, pack)
	return args.Error(0)
}

func (m *MockStore) SavePacks(ctx context.Context, packs []model.Pack, versionHash string) error {
	args := m.Called(ctx, packs, versionHash)
	return args.Error(0)
}

func (m *MockStore) GetPackByID(ctx context.Context, id string) (*model.Pack, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pack), args.Error(1)
}

func (m *MockStore) ListPacks(ctx context.Context) ([]model.Pack, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Pack), args.Error(1)
}

func (m *MockStore) DeletePack(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStore) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewPackService(t *testing.T) {
	mockStore := &MockStore{}
	service := NewPackService(mockStore)

	assert.NotNil(t, service)
	assert.Implements(t, (*PackService)(nil), service)
}

func TestPackService_CreatePacks(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		packs := []int64{250, 500, 1000}

		// Mock SavePacks to return success
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(nil)

		versionHash, err := service.CreatePacks(ctx, packs...)

		require.NoError(t, err)
		assert.NotEmpty(t, versionHash)
		assert.Len(t, versionHash, 16) // Hash is truncated to 16 characters

		mockStore.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		packs := []int64{250, 500}
		expectedError := errors.New("database error")

		// Mock SavePacks to return error
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(expectedError)

		versionHash, err := service.CreatePacks(ctx, packs...)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, versionHash)

		mockStore.AssertExpectations(t)
	})

	t.Run("empty packs", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		// Mock SavePacks with empty slice
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(nil)

		versionHash, err := service.CreatePacks(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, versionHash)

		mockStore.AssertExpectations(t)
	})

	t.Run("single pack", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		// Mock SavePacks
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(nil)

		versionHash, err := service.CreatePacks(ctx, 1000)

		require.NoError(t, err)
		assert.NotEmpty(t, versionHash)

		mockStore.AssertExpectations(t)
	})
}

func TestPackService_GetPacksByVersionHash(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedPacks := []model.Pack{
			{
				ID:          "pack-1",
				VersionHash: "abc123",
				TotalAmount: 250,
				PackItems: []model.PackItem{
					{ID: "item-1", PackID: "pack-1", Size: 250},
				},
			},
		}

		mockStore.On("GetPacksInvariantsByHash", ctx, "abc123").Return(expectedPacks, nil)

		result, err := service.GetPacksByVersionHash(ctx, "abc123")

		require.NoError(t, err)
		assert.Equal(t, expectedPacks, result)

		mockStore.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedError := store.ErrNotFound

		mockStore.On("GetPacksInvariantsByHash", ctx, "nonexistent").Return([]model.Pack{}, expectedError)

		result, err := service.GetPacksByVersionHash(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockStore.AssertExpectations(t)
	})
}

func TestPackService_GetPackByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockStore := &MockStore{}
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

		mockStore.On("GetPackByID", ctx, "pack-1").Return(expectedPack, nil)

		result, err := service.GetPackByID(ctx, "pack-1")

		require.NoError(t, err)
		assert.Equal(t, expectedPack, result)

		mockStore.AssertExpectations(t)
	})

	t.Run("pack not found", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		mockStore.On("GetPackByID", ctx, "nonexistent").Return(nil, store.ErrNotFound)

		result, err := service.GetPackByID(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, store.ErrNotFound, err)
		assert.Nil(t, result)

		mockStore.AssertExpectations(t)
	})
}

func TestPackService_ListPacks(t *testing.T) {
	t.Run("successful listing", func(t *testing.T) {
		mockStore := &MockStore{}
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

		mockStore.On("ListPacks", ctx).Return(expectedPacks, nil)

		result, err := service.ListPacks(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedPacks, result)

		mockStore.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		mockStore.On("ListPacks", ctx).Return([]model.Pack{}, nil)

		result, err := service.ListPacks(ctx)

		require.NoError(t, err)
		assert.Empty(t, result)

		mockStore.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedError := errors.New("database error")

		mockStore.On("ListPacks", ctx).Return([]model.Pack{}, expectedError)

		result, err := service.ListPacks(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockStore.AssertExpectations(t)
	})
}

func TestPackService_DeletePack(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		mockStore.On("DeletePack", ctx, "pack-1").Return(nil)

		err := service.DeletePack(ctx, "pack-1")

		require.NoError(t, err)

		mockStore.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)
		ctx := context.Background()

		expectedError := errors.New("database error")

		mockStore.On("DeletePack", ctx, "nonexistent").Return(expectedError)

		err := service.DeletePack(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockStore.AssertExpectations(t)
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

func TestPackService_ContextCancellation(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewPackService(mockStore)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// This would depend on how the store handles context cancellation
		// For now, we just test that the context is passed through
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(context.Canceled)

		_, err := service.CreatePacks(ctx, 250, 500)

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)

		mockStore.AssertExpectations(t)
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

func BenchmarkCreateInvariants(b *testing.B) {
	packs := []int64{250, 500, 1000, 2000, 5000}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateInvariants(packs...)
	}
}
