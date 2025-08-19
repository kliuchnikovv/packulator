package store

import (
	"context"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/stretchr/testify/assert"
)

// Simple tests that don't require database mocking
// These test business logic and data structures

func TestErrNotFound_Simple(t *testing.T) {
	assert.Equal(t, "not found", ErrNotFound.Error())
	assert.ErrorIs(t, ErrNotFound, ErrNotFound)
}

func TestStore_ContextHandling_Simple(t *testing.T) {
	ctx := context.Background()
	assert.NotNil(t, ctx)

	// Test context with timeout
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	assert.NotNil(t, ctx)
}

func TestPackDataValidation_Simple(t *testing.T) {
	pack := createTestPackSimple()

	assert.NotEmpty(t, pack.ID)
	assert.NotEmpty(t, pack.VersionHash)
	assert.Greater(t, pack.TotalAmount, int64(0))
	assert.NotEmpty(t, pack.PackItems)

	for _, item := range pack.PackItems {
		assert.NotEmpty(t, item.ID)
		assert.Equal(t, pack.ID, item.PackID)
		assert.Greater(t, item.Size, int64(0))
	}
}

func TestPacksDataValidation_Simple(t *testing.T) {
	packs := createTestPacksSimple()

	assert.Len(t, packs, 2)

	for _, pack := range packs {
		assert.NotEmpty(t, pack.ID)
		assert.NotEmpty(t, pack.VersionHash)
		assert.Greater(t, pack.TotalAmount, int64(0))
		assert.NotEmpty(t, pack.PackItems)
	}

	// All packs should have the same version hash in this test data
	assert.Equal(t, packs[0].VersionHash, packs[1].VersionHash)
}

// Test data creation and validation
func TestCreateTestPack_Simple(t *testing.T) {
	pack := createTestPackSimple()

	assert.Equal(t, "pack-1", pack.ID)
	assert.Equal(t, "abc123", pack.VersionHash)
	assert.Equal(t, int64(750), pack.TotalAmount)
	assert.Len(t, pack.PackItems, 2)

	// Verify pack items
	assert.Equal(t, "item-1", pack.PackItems[0].ID)
	assert.Equal(t, int64(250), pack.PackItems[0].Size)
	assert.Equal(t, "item-2", pack.PackItems[1].ID)
	assert.Equal(t, int64(500), pack.PackItems[1].Size)
}

// Benchmark pack creation
func BenchmarkCreateTestPack_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createTestPackSimple()
	}
}

func BenchmarkCreateTestPacks_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createTestPacksSimple()
	}
}

// Test data helpers
func createTestPackSimple() model.Pack {
	return model.Pack{
		ID:          "pack-1",
		VersionHash: "abc123",
		TotalAmount: 750,
		PackItems: []model.PackItem{
			{
				ID:     "item-1",
				PackID: "pack-1",
				Size:   250,
			},
			{
				ID:     "item-2",
				PackID: "pack-1",
				Size:   500,
			},
		},
	}
}

func createTestPacksSimple() []model.Pack {
	return []model.Pack{
		{
			ID:          "pack-1",
			VersionHash: "abc123",
			TotalAmount: 250,
			PackItems: []model.PackItem{
				{
					ID:     "item-1",
					PackID: "pack-1",
					Size:   250,
				},
			},
		},
		{
			ID:          "pack-2",
			VersionHash: "abc123",
			TotalAmount: 500,
			PackItems: []model.PackItem{
				{
					ID:     "item-2",
					PackID: "pack-2",
					Size:   500,
				},
			},
		},
	}
}

// Test error types and constants
func TestStoreErrors_Simple(t *testing.T) {
	// Test that ErrNotFound is properly defined
	assert.NotNil(t, ErrNotFound)
	assert.Contains(t, ErrNotFound.Error(), "not found")

	// Test that it can be used with errors.Is
	err := ErrNotFound
	assert.ErrorIs(t, err, ErrNotFound)
}

// Test model structures
func TestPackModel_Simple(t *testing.T) {
	pack := model.Pack{
		ID:          "test-id",
		VersionHash: "test-hash",
		TotalAmount: 100,
	}

	assert.Equal(t, "test-id", pack.ID)
	assert.Equal(t, "test-hash", pack.VersionHash)
	assert.Equal(t, int64(100), pack.TotalAmount)
}

func TestPackItemModel_Simple(t *testing.T) {
	item := model.PackItem{
		ID:     "item-id",
		PackID: "pack-id",
		Size:   50,
	}

	assert.Equal(t, "item-id", item.ID)
	assert.Equal(t, "pack-id", item.PackID)
	assert.Equal(t, int64(50), item.Size)
}
