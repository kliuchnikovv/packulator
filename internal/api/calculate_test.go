package api

import (
	"context"
	"errors"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/kliuchnikovv/packulator/internal/store"
	mock_store "github.com/kliuchnikovv/packulator/internal/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewPackagingService(t *testing.T) {
	mockStore := mock_store.NewMockStore(gomock.NewController(t))
	api := NewPackagingService(mockStore)

	assert.NotNil(t, api)
	assert.Equal(t, mockStore, api.store)
}

func TestPackagingService_Prefix(t *testing.T) {
	mockStore := mock_store.NewMockStore(gomock.NewController(t))
	api := NewPackagingService(mockStore)

	assert.Equal(t, "packaging", api.Prefix())
}

func TestPackagingService_NumberOfPackages(t *testing.T) {
	t.Run("successful calculation", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		// Test data
		amount := int64(1001)
		versionHash := "abc123"

		pack := model.Pack{
			ID:          "pack-1",
			VersionHash: versionHash,
			TotalAmount: 250,
			PackItems: []model.PackItem{
				{ID: "item-1", PackID: "pack-1", Size: 250},
				{ID: "item-2", PackID: "pack-2", Size: 500},
				{ID: "item-3", PackID: "pack-3", Size: 1000},
			},
		}

		// Mock request parameters
		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		// Mock store response
		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(&pack, nil)

		// Mock response
		response.On("OK", mock.AnythingOfType("map[int64]int64")).Return(nil)

		err := api.NumberOfPackages(ctx, request, response)

		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)

		// Verify the response data is a valid pack calculation
		responseData, ok := response.data.(map[int64]int64)
		require.True(t, ok)
		assert.NotEmpty(t, responseData)

		// Verify total amount is sufficient
		total := int64(0)
		for packSize, count := range responseData {
			total += packSize * count
		}
		assert.GreaterOrEqual(t, total, amount)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})

	t.Run("packs not found", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(1000)
		versionHash := "nonexistent"

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(nil, store.ErrNotFound)

		response.On("NotFound", "packs not found by hash: %s", mock.Anything).Return(store.ErrNotFound)

		err := api.NumberOfPackages(ctx, request, response)

		assert.Error(t, err)
		assert.Equal(t, store.ErrNotFound, err)
		assert.Equal(t, 404, response.statusCode)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})

	t.Run("store error", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(1000)
		versionHash := "abc123"
		expectedError := errors.New("database connection error")

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(nil, expectedError)

		response.On("InternalServerError", "failed to get packs: %s", mock.Anything).Return(expectedError)

		err := api.NumberOfPackages(ctx, request, response)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, 500, response.statusCode)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})

	t.Run("zero amount", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(0)
		versionHash := "abc123"

		pack := model.Pack{
			ID:          "pack-1",
			VersionHash: versionHash,
			TotalAmount: 250,
			PackItems: []model.PackItem{
				{ID: "item-1", PackID: "pack-1", Size: 250},
			},
		}

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(&pack, nil)

		response.On("OK", mock.AnythingOfType("map[int64]int64")).Return(nil)

		err := api.NumberOfPackages(ctx, request, response)

		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)

		// Response should be empty map for zero amount
		responseData, ok := response.data.(map[int64]int64)
		require.True(t, ok)
		assert.Empty(t, responseData)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})

	t.Run("large amount", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(50000)
		versionHash := "abc123"

		pack := model.Pack{
			ID:          "pack-1",
			VersionHash: versionHash,
			TotalAmount: 250,
			PackItems: []model.PackItem{
				{ID: "item-1", PackID: "pack-1", Size: 250},
				{ID: "item-2", PackID: "pack-2", Size: 5000},
			},
		}

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(&pack, nil)

		response.On("OK", mock.AnythingOfType("map[int64]int64")).Return(nil)

		err := api.NumberOfPackages(ctx, request, response)

		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)

		// Verify response has packs that cover the amount
		responseData, ok := response.data.(map[int64]int64)
		require.True(t, ok)
		assert.NotEmpty(t, responseData)

		total := int64(0)
		for packSize, count := range responseData {
			total += packSize * count
			assert.Greater(t, count, int64(0), "pack count should be positive")
		}
		assert.GreaterOrEqual(t, total, amount, "total should cover requested amount")

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})
}

func TestPackagingService_Routes(t *testing.T) {
	mockStore := mock_store.NewMockStore(gomock.NewController(t))
	api := NewPackagingService(mockStore)

	routes := api.Routers()

	// Check that routes are defined
	assert.NotEmpty(t, routes)

	// We can't easily test the exact route structure without more complex mocking
	// but we can verify the routes map is not empty
	assert.Greater(t, len(routes), 0, "should have at least one route defined")
}

func TestPackagingService_ContextHandling(t *testing.T) {
	t.Run("context with timeout", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(1000)
		versionHash := "abc123"

		pack := model.Pack{
			ID:          "pack-1",
			VersionHash: versionHash,
			TotalAmount: 1000,
			PackItems: []model.PackItem{
				{ID: "item-1", PackID: "pack-1", Size: 1000},
			},
		}

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(&pack, nil)

		response.On("OK", mock.AnythingOfType("map[int64]int64")).Return(nil)

		err := api.NumberOfPackages(ctx, request, response)

		require.NoError(t, err)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})

	t.Run("canceled context", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(1000)
		versionHash := "abc123"

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		// Mock store to return context canceled error
		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(nil, context.Canceled)
		response.On("InternalServerError", mock.Anything, mock.Anything).Return(context.Canceled)

		err := api.NumberOfPackages(ctx, request, response)

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})
}

func TestPackagingService_EdgeCases(t *testing.T) {
	t.Run("empty version hash", func(t *testing.T) {
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(1000)
		versionHash := ""

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(nil, store.ErrNotFound)
		response.On("NotFound", "packs not found by hash: %s", mock.Anything).Return(store.ErrNotFound)

		err := api.NumberOfPackages(ctx, request, response)

		assert.Error(t, err)
		assert.Equal(t, 404, response.statusCode)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})

	t.Run("negative amount", func(t *testing.T) {
		// Note: In a real application, this should be handled by input validation
		// but we can test the behavior if it somehow gets through
		mockStore := mock_store.NewMockStore(gomock.NewController(t))
		api := NewPackagingService(mockStore)
		ctx := context.Background()

		request := &MockRequest{}
		response := &MockResponse{}

		amount := int64(-100)
		versionHash := "abc123"

		pack := model.Pack{
			ID:          "pack-1",
			VersionHash: versionHash,
			TotalAmount: 250,
			PackItems: []model.PackItem{
				{ID: "item-1", PackID: "pack-1", Size: 250},
			},
		}

		request.On("Integer", "amount", mock.Anything).Return(amount)
		request.On("String", "packs_hash", mock.Anything).Return(versionHash)

		mockStore.EXPECT().GetPackByHash(gomock.Any(), versionHash).Return(&pack, nil)
		response.On("OK", mock.AnythingOfType("map[int64]int64")).Return(nil)

		err := api.NumberOfPackages(ctx, request, response)

		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)

		// Negative amount should result in empty map
		responseData, ok := response.data.(map[int64]int64)
		require.True(t, ok)
		assert.Empty(t, responseData)

		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})
}
