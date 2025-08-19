package api

import (
	"net/http"
	"testing"

	mock_store "github.com/kliuchnikovv/packulator/internal/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

func TestNewHealthAPI_Simple(t *testing.T) {
	store := mock_store.NewMockStore(gomock.NewController(t))
	api := NewHealthAPI(store)

	assert.NotNil(t, api)
	assert.Equal(t, store, api.store)
}

func TestHealthAPI_Prefix_Simple(t *testing.T) {
	api := &HealthAPI{}
	assert.Equal(t, "health", api.Prefix())
}

func TestHealthAPI_Middlewares_Simple(t *testing.T) {
	api := &HealthAPI{}
	middlewares := api.Middlewares()

	assert.Len(t, middlewares, 4)
}

func TestHealthAPI_Routers_Simple(t *testing.T) {
	api := &HealthAPI{}
	routes := api.Routers()

	assert.Len(t, routes, 1)
}

func TestHealthStatus_Struct(t *testing.T) {
	status := HealthStatus{
		Status:   "ok",
		Database: "ok",
		Version:  "1.0.0",
	}

	assert.Equal(t, "ok", status.Status)
	assert.Equal(t, "ok", status.Database)
	assert.Equal(t, "1.0.0", status.Version)
}

type mockResponse struct {
	mock.Mock
	statusCode int
	data       interface{}
}

func (m *mockResponse) OK(data interface{}) error {
	args := m.Called(data)
	m.statusCode = 200
	m.data = data
	return args.Error(0)
}

func (m *mockResponse) Created() error {
	args := m.Called()
	m.statusCode = 201
	m.data = nil
	return args.Error(0)
}

func (m *mockResponse) NoContent() error {
	args := m.Called()
	m.statusCode = 204
	m.data = nil
	return args.Error(0)
}

func (m *mockResponse) BadRequest(format string, args ...interface{}) error {
	callArgs := m.Called(format, args)
	m.statusCode = 400
	return callArgs.Error(0)
}

func (m *mockResponse) InternalServerError(format string, args ...interface{}) error {
	callArgs := m.Called(format, args)
	m.statusCode = 500
	return callArgs.Error(0)
}

func (m *mockResponse) NotFound(format string, args ...interface{}) error {
	callArgs := m.Called(format, args)
	m.statusCode = 404
	return callArgs.Error(0)
}

func (m *mockResponse) Object(code int, payload interface{}) error {
	args := m.Called(code, payload)
	m.statusCode = code
	m.data = payload
	return args.Error(0)
}

func (m *mockResponse) Error(code int, err error) error {
	args := m.Called(code, err)
	m.statusCode = code
	return args.Error(0)
}

func (m *mockResponse) Errorf(code int, format string, args ...interface{}) error {
	callArgs := m.Called(code, format, args)
	m.statusCode = code
	return callArgs.Error(0)
}

func (m *mockResponse) Forbidden(format string, args ...interface{}) error {
	callArgs := m.Called(format, args)
	m.statusCode = 403
	return callArgs.Error(0)
}

func (m *mockResponse) MethodNotAllowed(format string, args ...interface{}) error {
	callArgs := m.Called(format, args)
	m.statusCode = 405
	return callArgs.Error(0)
}

func (m *mockResponse) ResponseWriter() http.ResponseWriter {
	args := m.Called()
	return args.Get(0).(http.ResponseWriter)
}

func TestNewHealthAPI(t *testing.T) {
	store := mock_store.NewMockStore(gomock.NewController(t))
	api := NewHealthAPI(store)

	assert.NotNil(t, api)
	assert.Equal(t, store, api.store)
}

func TestHealthAPI_Prefix(t *testing.T) {
	api := &HealthAPI{}
	assert.Equal(t, "health", api.Prefix())
}

func TestHealthAPI_Middlewares(t *testing.T) {
	api := &HealthAPI{}
	middlewares := api.Middlewares()

	assert.Len(t, middlewares, 4)
}

func TestHealthAPI_Routers(t *testing.T) {
	api := &HealthAPI{}
	routes := api.Routers()

	assert.Len(t, routes, 1)
}
