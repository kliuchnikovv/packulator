package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/kliuchnikovv/engi/definition/parameter/placing"
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

// Mock response for testing
type MockResponse struct {
	mock.Mock
	statusCode int
	data       interface{}
}

func (m *MockResponse) OK(data interface{}) error {
	args := m.Called(data)
	m.statusCode = 200
	m.data = data
	return args.Error(0)
}

func (m *MockResponse) Created() error {
	args := m.Called()
	m.statusCode = 201
	m.data = nil
	return args.Error(0)
}

func (m *MockResponse) NoContent() error {
	args := m.Called()
	m.statusCode = 204
	m.data = nil
	return args.Error(0)
}

func (m *MockResponse) BadRequest(format string, args ...interface{}) error {
	mockArgs := []interface{}{format}
	mockArgs = append(mockArgs, args...)
	callArgs := m.Called(mockArgs...)
	m.statusCode = 400
	return callArgs.Error(0)
}

func (m *MockResponse) NotFound(format string, args ...interface{}) error {
	mockArgs := []interface{}{format}
	mockArgs = append(mockArgs, args...)
	callArgs := m.Called(mockArgs...)
	m.statusCode = 404
	return callArgs.Error(0)
}

func (m *MockResponse) InternalServerError(format string, args ...interface{}) error {
	mockArgs := []interface{}{format}
	mockArgs = append(mockArgs, args...)
	callArgs := m.Called(mockArgs...)
	m.statusCode = 500
	return callArgs.Error(0)
}

func (m *MockResponse) ResponseWriter() http.ResponseWriter {
	args := m.Called()
	return args.Get(0).(http.ResponseWriter)
}

func (m *MockResponse) Object(code int, payload interface{}) error {
	args := m.Called(code, payload)
	m.statusCode = code
	m.data = payload
	return args.Error(0)
}

func (m *MockResponse) WithoutContent(code int) error {
	args := m.Called(code)
	m.statusCode = code
	m.data = nil
	return args.Error(0)
}

func (m *MockResponse) Error(code int, err error) error {
	args := m.Called(code, err)
	m.statusCode = code
	return args.Error(0)
}

func (m *MockResponse) Errorf(code int, format string, args ...interface{}) error {
	mockArgs := []interface{}{code, format}
	mockArgs = append(mockArgs, args...)
	callArgs := m.Called(mockArgs...)
	m.statusCode = code
	return callArgs.Error(0)
}

func (m *MockResponse) Forbidden(format string, args ...interface{}) error {
	mockArgs := []interface{}{format}
	mockArgs = append(mockArgs, args...)
	callArgs := m.Called(mockArgs...)
	m.statusCode = 403
	return callArgs.Error(0)
}

func (m *MockResponse) MethodNotAllowed(format string, args ...interface{}) error {
	mockArgs := []interface{}{format}
	mockArgs = append(mockArgs, args...)
	callArgs := m.Called(mockArgs...)
	m.statusCode = 405
	return callArgs.Error(0)
}

// Mock request for testing
type MockRequest struct {
	mock.Mock
	body interface{}
}

func (m *MockRequest) Body() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockRequest) String(key string, paramPlacing placing.Placing) string {
	args := m.Called(key, paramPlacing)
	return args.String(0)
}

func (m *MockRequest) Integer(key string, paramPlacing placing.Placing) int64 {
	args := m.Called(key, paramPlacing)
	return args.Get(0).(int64)
}

func (m *MockRequest) Bool(key string, paramPlacing placing.Placing) bool {
	args := m.Called(key, paramPlacing)
	return args.Bool(0)
}

func (m *MockRequest) Float(key string, paramPlacing placing.Placing) float64 {
	args := m.Called(key, paramPlacing)
	return args.Get(0).(float64)
}

func (m *MockRequest) Headers() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *MockRequest) Parameters() map[placing.Placing]map[string]string {
	args := m.Called()
	return args.Get(0).(map[placing.Placing]map[string]string)
}

func (m *MockRequest) GetParameter(value string, place placing.Placing) string {
	args := m.Called(value, place)
	return args.String(0)
}

func (m *MockRequest) GetRequest() *http.Request {
	args := m.Called()
	return args.Get(0).(*http.Request)
}

func (m *MockRequest) Time(key string, layout string, paramPlacing placing.Placing) time.Time {
	args := m.Called(key, layout, paramPlacing)
	return args.Get(0).(time.Time)
}

func TestNewPacksAPI(t *testing.T) {
	mockStore := &MockStore{}
	api := NewPacksAPI(mockStore)
	
	assert.NotNil(t, api)
	assert.NotNil(t, api.packService)
}

func TestPacksAPI_Prefix(t *testing.T) {
	mockStore := &MockStore{}
	api := NewPacksAPI(mockStore)
	
	assert.Equal(t, "packs", api.Prefix())
}

func TestPacksAPI_CreatePacks(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		requestBody := &model.CreatePacksRequest{
			Packs: []int64{250, 500, 1000},
		}
		
		// Mock request body
		request.On("Body").Return(requestBody)
		
		// Mock store SavePacks
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(nil)
		
		// Mock response
		response.On("OK", mock.AnythingOfType("model.CreatePacksResponse")).Return(nil)
		
		err := api.CreatePacks(ctx, request, response)
		
		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)
		
		// Verify response data
		responseData, ok := response.data.(model.CreatePacksResponse)
		require.True(t, ok)
		assert.NotEmpty(t, responseData.VersionHash)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
	
	t.Run("empty packs", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		requestBody := &model.CreatePacksRequest{
			Packs: []int64{},
		}
		
		request.On("Body").Return(requestBody)
		response.On("InternalServerError", mock.Anything).Return(errors.New("packs can't be empty"))
		
		err := api.CreatePacks(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, 500, response.statusCode)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
	})
	
	t.Run("service error", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		requestBody := &model.CreatePacksRequest{
			Packs: []int64{250, 500},
		}
		
		request.On("Body").Return(requestBody)
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(errors.New("database error"))
		response.On("InternalServerError", mock.Anything, mock.Anything).Return(errors.New("can't create packs"))
		
		err := api.CreatePacks(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, 500, response.statusCode)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}

func TestPacksAPI_ListPacks(t *testing.T) {
	t.Run("successful listing", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
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
		response.On("OK", expectedPacks).Return(nil)
		
		err := api.ListPacks(ctx, request, response)
		
		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)
		assert.Equal(t, expectedPacks, response.data)
		
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
	
	t.Run("service error", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		expectedError := errors.New("database error")
		mockStore.On("ListPacks", ctx).Return([]model.Pack{}, expectedError)
		response.On("InternalServerError", mock.Anything, mock.Anything).Return(expectedError)
		
		err := api.ListPacks(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, 500, response.statusCode)
		
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}

func TestPacksAPI_GetPackByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		expectedPack := &model.Pack{
			ID:          "pack-1",
			VersionHash: "abc123",
			TotalAmount: 250,
		}
		
		request.On("String", "id", mock.Anything).Return("pack-1")
		mockStore.On("GetPackByID", ctx, "pack-1").Return(expectedPack, nil)
		response.On("OK", expectedPack).Return(nil)
		
		err := api.GetPackByID(ctx, request, response)
		
		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)
		assert.Equal(t, expectedPack, response.data)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
	
	t.Run("pack not found", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		request.On("String", "id", mock.Anything).Return("nonexistent")
		mockStore.On("GetPackByID", ctx, "nonexistent").Return(nil, store.ErrNotFound)
		response.On("InternalServerError", mock.Anything, mock.Anything, mock.Anything).Return(store.ErrNotFound)
		
		err := api.GetPackByID(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, 500, response.statusCode)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}

func TestPacksAPI_GetPacksByVersionHash(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		expectedPacks := []model.Pack{
			{
				ID:          "pack-1",
				VersionHash: "abc123",
				TotalAmount: 250,
			},
		}
		
		request.On("String", "hash", mock.Anything).Return("abc123")
		mockStore.On("GetPacksInvariantsByHash", ctx, "abc123").Return(expectedPacks, nil)
		response.On("OK", expectedPacks).Return(nil)
		
		err := api.GetPacksByVersionHash(ctx, request, response)
		
		require.NoError(t, err)
		assert.Equal(t, 200, response.statusCode)
		assert.Equal(t, expectedPacks, response.data)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
	
	t.Run("hash not found", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		request.On("String", "hash", mock.Anything).Return("nonexistent")
		mockStore.On("GetPacksInvariantsByHash", ctx, "nonexistent").Return([]model.Pack{}, store.ErrNotFound)
		response.On("InternalServerError", mock.Anything, mock.Anything, mock.Anything).Return(store.ErrNotFound)
		
		err := api.GetPacksByVersionHash(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, 500, response.statusCode)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}

func TestPacksAPI_DeletePack(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		request.On("String", "id", mock.Anything).Return("pack-1")
		mockStore.On("DeletePack", ctx, "pack-1").Return(nil)
		response.On("NoContent").Return(nil)
		
		err := api.DeletePack(ctx, request, response)
		
		require.NoError(t, err)
		assert.Equal(t, 204, response.statusCode)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
	
	t.Run("service error", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		ctx := context.Background()
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		expectedError := errors.New("database error")
		request.On("String", "id", mock.Anything).Return("pack-1")
		mockStore.On("DeletePack", ctx, "pack-1").Return(expectedError)
		response.On("InternalServerError", mock.Anything, mock.Anything).Return(expectedError)
		
		err := api.DeletePack(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, 500, response.statusCode)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}

func TestPacksAPI_Routes(t *testing.T) {
	mockStore := &MockStore{}
	api := NewPacksAPI(mockStore)
	
	routes := api.Routers()
	
	// Check that routes are defined
	assert.NotEmpty(t, routes)
	
	// We can't easily test the exact route structure without more complex mocking
	// but we can verify the routes map is not empty
	assert.Greater(t, len(routes), 0, "should have at least one route defined")
}

func TestPacksAPI_ContextHandling(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		mockStore := &MockStore{}
		api := NewPacksAPI(mockStore)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		request := &MockRequest{}
		response := &MockResponse{}
		
		// Mock the behavior for cancelled context
		request.On("Body").Return(&model.CreatePacksRequest{Packs: []int64{250}})
		mockStore.On("SavePacks", ctx, mock.AnythingOfType("[]model.Pack"), mock.AnythingOfType("string")).Return(context.Canceled)
		response.On("InternalServerError", mock.Anything, mock.Anything).Return(context.Canceled)
		
		err := api.CreatePacks(ctx, request, response)
		
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		
		request.AssertExpectations(t)
		response.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}