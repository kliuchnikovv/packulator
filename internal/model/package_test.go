package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPack_GetPacks(t *testing.T) {
	tests := []struct {
		name      string
		pack      Pack
		expected  []int64
	}{
		{
			name: "empty pack items",
			pack: Pack{
				ID:          "pack-1",
				TotalAmount: 0,
				PackItems:   []PackItem{},
			},
			expected: []int64{},
		},
		{
			name: "single pack item",
			pack: Pack{
				ID:          "pack-1",
				TotalAmount: 250,
				PackItems: []PackItem{
					{
						ID:     "item-1",
						PackID: "pack-1",
						Size:   250,
					},
				},
			},
			expected: []int64{250},
		},
		{
			name: "multiple pack items",
			pack: Pack{
				ID:          "pack-1",
				TotalAmount: 750,
				PackItems: []PackItem{
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
			},
			expected: []int64{250, 500},
		},
		{
			name: "pack items with various sizes",
			pack: Pack{
				ID:          "pack-1",
				TotalAmount: 7750,
				PackItems: []PackItem{
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
					{
						ID:     "item-3",
						PackID: "pack-1",
						Size:   1000,
					},
					{
						ID:     "item-4",
						PackID: "pack-1",
						Size:   2000,
					},
					{
						ID:     "item-5",
						PackID: "pack-1",
						Size:   4000,
					},
				},
			},
			expected: []int64{250, 500, 1000, 2000, 4000},
		},
		{
			name: "pack items with duplicate sizes",
			pack: Pack{
				ID:          "pack-1",
				TotalAmount: 500,
				PackItems: []PackItem{
					{
						ID:     "item-1",
						PackID: "pack-1",
						Size:   250,
					},
					{
						ID:     "item-2",
						PackID: "pack-1",
						Size:   250,
					},
				},
			},
			expected: []int64{250, 250},
		},
		{
			name: "pack items with zero size",
			pack: Pack{
				ID:          "pack-1",
				TotalAmount: 250,
				PackItems: []PackItem{
					{
						ID:     "item-1",
						PackID: "pack-1",
						Size:   0,
					},
					{
						ID:     "item-2",
						PackID: "pack-1",
						Size:   250,
					},
				},
			},
			expected: []int64{0, 250},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pack.GetPacks()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPack_GetPacks_PreservesOrder(t *testing.T) {
	pack := Pack{
		ID:          "pack-1",
		TotalAmount: 1750,
		PackItems: []PackItem{
			{ID: "item-1", PackID: "pack-1", Size: 1000},
			{ID: "item-2", PackID: "pack-1", Size: 250},
			{ID: "item-3", PackID: "pack-1", Size: 500},
		},
	}

	result := pack.GetPacks()
	expected := []int64{1000, 250, 500}

	assert.Equal(t, expected, result, "GetPacks should preserve the order of PackItems")
}

func TestPack_GetPacks_EmptySliceNotNil(t *testing.T) {
	pack := Pack{
		ID:          "pack-1",
		TotalAmount: 0,
		PackItems:   []PackItem{},
	}

	result := pack.GetPacks()
	
	assert.NotNil(t, result, "GetPacks should return an empty slice, not nil")
	assert.Len(t, result, 0, "GetPacks should return an empty slice")
	assert.IsType(t, []int64{}, result, "GetPacks should return []int64 type")
}

func TestPack_GetPacks_WithNilPackItems(t *testing.T) {
	pack := Pack{
		ID:          "pack-1",
		TotalAmount: 0,
		PackItems:   nil,
	}

	result := pack.GetPacks()
	
	assert.NotNil(t, result, "GetPacks should return an empty slice, not nil")
	assert.Len(t, result, 0, "GetPacks should return an empty slice")
}

func TestPackItem_Fields(t *testing.T) {
	item := PackItem{
		ID:     "item-123",
		PackID: "pack-456",
		Size:   1000,
	}

	assert.Equal(t, "item-123", item.ID)
	assert.Equal(t, "pack-456", item.PackID)
	assert.Equal(t, int64(1000), item.Size)
}

func TestPack_Fields(t *testing.T) {
	pack := Pack{
		ID:          "pack-123",
		VersionHash: "abc123def456",
		TotalAmount: 1500,
		PackItems: []PackItem{
			{ID: "item-1", PackID: "pack-123", Size: 500},
			{ID: "item-2", PackID: "pack-123", Size: 1000},
		},
	}

	assert.Equal(t, "pack-123", pack.ID)
	assert.Equal(t, "abc123def456", pack.VersionHash)
	assert.Equal(t, int64(1500), pack.TotalAmount)
	assert.Len(t, pack.PackItems, 2)
}

func TestPackItem_ZeroValues(t *testing.T) {
	var item PackItem
	
	assert.Equal(t, "", item.ID)
	assert.Equal(t, "", item.PackID)
	assert.Equal(t, int64(0), item.Size)
}

func TestPack_ZeroValues(t *testing.T) {
	var pack Pack
	
	assert.Equal(t, "", pack.ID)
	assert.Equal(t, "", pack.VersionHash)
	assert.Equal(t, int64(0), pack.TotalAmount)
	assert.Nil(t, pack.PackItems)
}