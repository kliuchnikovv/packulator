package service

import (
	"testing"

	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvariants(t *testing.T) {
	cases := []struct {
		name     string
		input    []int64
		expected []model.Pack
	}{
		{
			name:     "empty",
			input:    []int64{},
			expected: []model.Pack{},
		},
		{
			name:  "single",
			input: []int64{1},
			expected: []model.Pack{
				{
					TotalAmount: 1,
					PackItems: []model.PackItem{{
						Size: 1,
					}},
				},
			},
		},
		{
			name:  "multiple",
			input: []int64{1, 2, 3},
			expected: []model.Pack{
				{
					TotalAmount: 1,
					PackItems: []model.PackItem{{
						Size: 1,
					}},
				},
				{
					TotalAmount: 2,
					PackItems: []model.PackItem{{
						Size: 2,
					}},
				},
				{
					TotalAmount: 3,
					PackItems: []model.PackItem{{
						Size: 3,
					}},
				},
				{
					TotalAmount: 3,
					PackItems: []model.PackItem{{
						Size: 1,
					}, {
						Size: 2,
					}},
				},
				{
					TotalAmount: 4,
					PackItems: []model.PackItem{{
						Size: 1,
					}, {
						Size: 3,
					}},
				},
				{
					TotalAmount: 5,
					PackItems: []model.PackItem{{
						Size: 2,
					}, {
						Size: 3,
					}},
				},
				{
					TotalAmount: 6,
					PackItems: []model.PackItem{{
						Size: 1,
					}, {
						Size: 2,
					}, {
						Size: 3,
					}},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateInvariants(tt.input...)
			
			// Compare lengths first
			assert.Equal(t, len(tt.expected), len(result), "length should match")
			
			// Compare each pack, ignoring IDs and UUIDs
			for i, expectedPack := range tt.expected {
				if i >= len(result) {
					break
				}
				actualPack := result[i]
				
				// Compare TotalAmount
				assert.Equal(t, expectedPack.TotalAmount, actualPack.TotalAmount, "pack %d total amount", i)
				
				// Compare PackItems length
				assert.Equal(t, len(expectedPack.PackItems), len(actualPack.PackItems), "pack %d items length", i)
				
				// Compare each PackItem size (ignoring IDs)
				for j, expectedItem := range expectedPack.PackItems {
					if j >= len(actualPack.PackItems) {
						break
					}
					assert.Equal(t, expectedItem.Size, actualPack.PackItems[j].Size, "pack %d item %d size", i, j)
				}
			}
		})
	}
}

func TestCreateInvariantsEdgeCases(t *testing.T) {
	t.Run("duplicate pack sizes", func(t *testing.T) {
		result := CreateInvariants(1, 1, 2)
		
		// Should handle duplicates correctly
		assert.Greater(t, len(result), 0)
		
		// Check that we have some basic invariants
		foundOne := false
		foundTwo := false
		for _, pack := range result {
			if pack.TotalAmount == 1 && len(pack.PackItems) == 1 {
				foundOne = true
			}
			if pack.TotalAmount == 2 && len(pack.PackItems) == 1 {
				foundTwo = true
			}
		}
		assert.True(t, foundOne, "should have single pack with size 1")
		assert.True(t, foundTwo, "should have single pack with size 2")
	})

	t.Run("unsorted input", func(t *testing.T) {
		result := CreateInvariants(3, 1, 2)
		
		// Should still produce correct invariants regardless of input order
		assert.Greater(t, len(result), 0)
		
		// Check for expected invariants
		foundSmallest := false
		foundLargest := false
		for _, pack := range result {
			if pack.TotalAmount == 1 && len(pack.PackItems) == 1 {
				foundSmallest = true
			}
			if pack.TotalAmount == 3 && len(pack.PackItems) == 1 {
				foundLargest = true
			}
		}
		assert.True(t, foundSmallest, "should have smallest pack")
		assert.True(t, foundLargest, "should have largest pack")
	})
}

func TestSortInvariants(t *testing.T) {
	packs := []model.Pack{
		{TotalAmount: 3, PackItems: []model.PackItem{{Size: 1}, {Size: 2}}},
		{TotalAmount: 1, PackItems: []model.PackItem{{Size: 1}}},
		{TotalAmount: 3, PackItems: []model.PackItem{{Size: 3}}},
		{TotalAmount: 2, PackItems: []model.PackItem{{Size: 2}}},
	}

	t.Run("ascending order", func(t *testing.T) {
		testPacks := make([]model.Pack, len(packs))
		copy(testPacks, packs)
		
		sortInvariants(testPacks, true)
		
		// Check ascending order by TotalAmount
		assert.Equal(t, int64(1), testPacks[0].TotalAmount)
		assert.Equal(t, int64(2), testPacks[1].TotalAmount)
		
		// For same TotalAmount, shorter PackItems should come first
		for i := 2; i < len(testPacks)-1; i++ {
			if testPacks[i].TotalAmount == testPacks[i+1].TotalAmount {
				assert.LessOrEqual(t, len(testPacks[i].PackItems), len(testPacks[i+1].PackItems))
			}
		}
	})

	t.Run("descending order", func(t *testing.T) {
		testPacks := make([]model.Pack, len(packs))
		copy(testPacks, packs)
		
		sortInvariants(testPacks, false)
		
		// Check descending order by TotalAmount
		for i := 0; i < len(testPacks)-1; i++ {
			assert.GreaterOrEqual(t, testPacks[i].TotalAmount, testPacks[i+1].TotalAmount)
		}
	})
}
