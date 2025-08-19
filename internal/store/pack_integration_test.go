//go:build integration
// +build integration

package store

import (
	"context"
	"os"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/config"
	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
)

// Integration tests require a real database
// Run with: go test -tags=integration ./internal/store

func getTestDatabaseConfig() *config.DatabaseConfig {
	// Use environment variables or defaults for test database
	return &config.DatabaseConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		Database: getEnvOrDefault("TEST_DB_NAME", "packulator_test"),
		SSLMode:  getEnvOrDefault("TEST_DB_SSL", "disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupTestStore(t *testing.T) Store {
	cfg := getTestDatabaseConfig()
	store, err := NewStore(postgres.Open(cfg.DSN()))
	require.NoError(t, err, "Failed to create test store")
	return store
}

func createTestPackForIntegration() *model.Pack {
	return &model.Pack{
		ID:          "test-pack-integration",
		VersionHash: "test-hash-123",
		TotalAmount: 750,
		PackItems: []model.PackItem{
			{
				ID:     "test-item-1",
				PackID: "test-pack-integration",
				Size:   250,
			},
			{
				ID:     "test-item-2",
				PackID: "test-pack-integration",
				Size:   500,
			},
		},
	}
}

func TestStoreIntegration_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	err := store.HealthCheck(ctx)
	assert.NoError(t, err, "Health check should pass with valid database connection")
}

func TestStoreIntegration_SaveAndGetPack(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	// Create test pack
	pack := createTestPackForIntegration()

	// Save pack
	err := store.SavePack(ctx, pack)
	require.NoError(t, err, "Should save pack successfully")

	// Get pack by ID
	retrievedPack, err := store.GetPackByID(ctx, pack.ID)
	require.NoError(t, err, "Should retrieve pack successfully")

	assert.Equal(t, pack.ID, retrievedPack.ID)
	assert.Equal(t, pack.VersionHash, retrievedPack.VersionHash)
	assert.Equal(t, pack.TotalAmount, retrievedPack.TotalAmount)
	assert.Len(t, retrievedPack.PackItems, 2)

	// Cleanup
	err = store.DeletePack(ctx, pack.ID)
	require.NoError(t, err, "Should delete pack successfully")
}

func TestStoreIntegration_SavePacks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	versionHash := "batch-test-hash"
	pack1 := model.Pack{
		ID:          "batch-pack-1",
		VersionHash: versionHash,
		TotalAmount: 250,
		PackItems: []model.PackItem{
			{
				ID:     "batch-item-1",
				PackID: "batch-pack-1",
				Size:   250,
			},
		},
	}

	pack2 := model.Pack{
		ID:          "batch-pack-2",
		VersionHash: versionHash,
		TotalAmount: 500,
		PackItems: []model.PackItem{
			{
				ID:     "batch-item-2",
				PackID: "batch-pack-2",
				Size:   500,
			},
		},
	}

	// Save packs in batch
	err := store.SavePacks(ctx, pack1, pack2)
	require.NoError(t, err, "Should save packs in batch successfully")

	// Get pack by version hash (should find the first one)
	retrievedPack, err := store.GetPackByHash(ctx, versionHash)
	require.NoError(t, err, "Should retrieve pack by hash successfully")

	// Verify pack has correct version hash
	assert.Equal(t, versionHash, retrievedPack.VersionHash)
	assert.NotEmpty(t, retrievedPack.PackItems)

	// Cleanup both packs
	err = store.DeletePack(ctx, pack1.ID)
	require.NoError(t, err, "Should delete pack1 successfully")
	err = store.DeletePack(ctx, pack2.ID)
	require.NoError(t, err, "Should delete pack2 successfully")
}

func TestStoreIntegration_ListPacks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	// Save a test pack
	pack := createTestPackForIntegration()
	pack.ID = "list-test-pack"

	err := store.SavePack(ctx, pack)
	require.NoError(t, err, "Should save pack successfully")

	// List all packs
	packs, err := store.ListPacks(ctx)
	require.NoError(t, err, "Should list packs successfully")

	// Should have at least our test pack
	assert.GreaterOrEqual(t, len(packs), 1)

	// Find our test pack in the list
	found := false
	for _, p := range packs {
		if p.ID == pack.ID {
			found = true
			assert.Equal(t, pack.VersionHash, p.VersionHash)
			assert.Equal(t, pack.TotalAmount, p.TotalAmount)
			break
		}
	}
	assert.True(t, found, "Should find our test pack in the list")

	// Cleanup
	err = store.DeletePack(ctx, pack.ID)
	require.NoError(t, err, "Should delete pack successfully")
}

func TestStoreIntegration_GetPackByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	// Try to get non-existent pack
	pack, err := store.GetPackByID(ctx, "non-existent-pack")

	assert.Error(t, err, "Should return error for non-existent pack")
	assert.ErrorIs(t, err, ErrNotFound, "Should return ErrNotFound specifically")
	assert.Nil(t, pack, "Pack should be nil")
}

func TestStoreIntegration_GetPackByHash_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	// Try to get pack with non-existent hash
	pack, err := store.GetPackByHash(ctx, "non-existent-hash")

	assert.Error(t, err, "Should return error for non-existent hash")
	assert.ErrorIs(t, err, ErrNotFound, "Should return ErrNotFound specifically")
	assert.Nil(t, pack, "Pack should be nil")
}

func TestStoreIntegration_DeletePack_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestStore(t)
	ctx := context.Background()

	// Try to delete non-existent pack
	err := store.DeletePack(ctx, "non-existent-pack")

	// GORM delete doesn't return error for non-existent records
	assert.NoError(t, err, "Delete should not return error even for non-existent pack")
}
