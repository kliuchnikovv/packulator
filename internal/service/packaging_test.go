package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPackagingService(t *testing.T) {
	service := NewPackagingService()
	assert.NotNil(t, service)
	assert.Implements(t, (*PackagingService)(nil), service)
}

func TestNumberOfPacks(t *testing.T) {
	var (
		s          = NewPackagingService()
		ctx        = context.TODO()
		invariants = CreateInvariants(250, 500, 1000, 2000, 5000)
	)

	cases := []struct {
		amount   int64
		expected map[int64]int64
	}{
		{
			amount: 1,
			expected: map[int64]int64{
				250: 1,
			},
		},
		{
			amount: 250,
			expected: map[int64]int64{
				250: 1,
			},
		},
		{
			amount: 251,
			expected: map[int64]int64{
				500: 1,
			},
		},
		{
			amount: 501,
			expected: map[int64]int64{
				500: 1,
				250: 1,
			},
		},
		{
			amount: 12001,
			expected: map[int64]int64{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
	}

	for _, tc := range cases {
		t.Run(strconv.FormatInt(tc.amount, 10), func(t *testing.T) {
			actual, err := s.NumberOfPacks(ctx, tc.amount, invariants)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestNumberOfPacks_EdgeCases(t *testing.T) {
	s := NewPackagingService()
	ctx := context.Background()

	t.Run("zero amount", func(t *testing.T) {
		invariants := CreateInvariants(250, 500)
		result, err := s.NumberOfPacks(ctx, 0, invariants)

		require.NoError(t, err)
		assert.Empty(t, result, "zero amount should result in empty pack map")
	})

	t.Run("negative amount", func(t *testing.T) {
		invariants := CreateInvariants(250, 500)
		result, err := s.NumberOfPacks(ctx, -100, invariants)

		require.NoError(t, err)
		assert.Empty(t, result, "negative amount should result in empty pack map")
	})

	t.Run("empty invariants", func(t *testing.T) {
		invariants := []model.Pack{}
		result, err := s.NumberOfPacks(ctx, 1000, invariants)

		require.NoError(t, err)
		assert.Empty(t, result, "empty invariants should result in empty pack map")
	})

	t.Run("single pack size", func(t *testing.T) {
		invariants := CreateInvariants(100)
		result, err := s.NumberOfPacks(ctx, 250, invariants)

		require.NoError(t, err)
		assert.Equal(t, map[int64]int64{100: 3}, result, "should use 3 packs of size 100")
	})

	t.Run("exact match", func(t *testing.T) {
		invariants := CreateInvariants(250, 500, 1000)
		result, err := s.NumberOfPacks(ctx, 1000, invariants)

		require.NoError(t, err)
		assert.Equal(t, map[int64]int64{1000: 1}, result, "should use exactly one 1000-size pack")
	})

	t.Run("very small amount", func(t *testing.T) {
		invariants := CreateInvariants(250, 500, 1000)
		result, err := s.NumberOfPacks(ctx, 1, invariants)

		require.NoError(t, err)
		assert.Equal(t, map[int64]int64{250: 1}, result, "should use smallest available pack")
	})

	t.Run("large amount", func(t *testing.T) {
		invariants := CreateInvariants(250, 500, 1000, 2000, 5000)
		result, err := s.NumberOfPacks(ctx, 50000, invariants)

		require.NoError(t, err)
		assert.NotEmpty(t, result)

		// Verify total amount
		total := int64(0)
		for packSize, count := range result {
			total += packSize * count
		}
		assert.GreaterOrEqual(t, total, int64(50000), "total should cover the requested amount")
	})
}

func TestNumberOfPacks_DifferentPackSizes(t *testing.T) {
	s := NewPackagingService()
	ctx := context.Background()

	t.Run("custom pack sizes - small", func(t *testing.T) {
		invariants := CreateInvariants(10, 25, 50)
		result, err := s.NumberOfPacks(ctx, 75, invariants)

		require.NoError(t, err)
		assert.NotEmpty(t, result)

		// Verify we can achieve the target amount
		total := int64(0)
		for packSize, count := range result {
			total += packSize * count
		}
		assert.GreaterOrEqual(t, total, int64(75))
	})

	t.Run("custom pack sizes - prime numbers", func(t *testing.T) {
		invariants := CreateInvariants(7, 11, 13, 17)
		result, err := s.NumberOfPacks(ctx, 100, invariants)

		require.NoError(t, err)
		assert.NotEmpty(t, result)

		// Verify we can achieve the target amount
		total := int64(0)
		for packSize, count := range result {
			total += packSize * count
		}
		assert.GreaterOrEqual(t, total, int64(100))
	})

	t.Run("duplicate pack sizes", func(t *testing.T) {
		invariants := CreateInvariants(100, 100, 200)
		result, err := s.NumberOfPacks(ctx, 350, invariants)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestNumberOfPacks_ContextHandling(t *testing.T) {
	s := NewPackagingService()
	invariants := CreateInvariants(250, 500, 1000)

	t.Run("context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		result, err := s.NumberOfPacks(ctx, 1000, invariants)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// The current implementation doesn't check context cancellation
		// but we test the interface
		result, err := s.NumberOfPacks(ctx, 1000, invariants)

		// Current implementation doesn't handle context cancellation
		require.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestNumberOfPacks_IntegrationWithCreateInvariants(t *testing.T) {
	s := NewPackagingService()
	ctx := context.Background()

	t.Run("end-to-end with different pack configurations", func(t *testing.T) {
		testCases := []struct {
			name      string
			packSizes []int64
			amount    int64
			expected  map[int64]int64
		}{
			// {
			// 	name:      "standard packs",
			// 	packSizes: []int64{250, 500, 1000, 2000, 5000},
			// 	amount:    3750,
			// 	expected:  map[int64]int64{2000: 1, 1000: 1, 500: 1, 250: 1},
			// },
			{
				name:      "small packs",
				packSizes: []int64{5, 10, 25},
				amount:    100,
				expected:  map[int64]int64{25: 4},
			},
			{
				name:      "large packs",
				packSizes: []int64{10000, 50000, 100000},
				amount:    150000,
				expected:  map[int64]int64{100000: 1, 50000: 1},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				invariants := CreateInvariants(tc.packSizes...)
				result, err := s.NumberOfPacks(ctx, tc.amount, invariants)

				require.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Equal(t, tc.expected, result)
			})
		}
	})
}

// Benchmark tests
func BenchmarkNumberOfPacks_Small(b *testing.B) {
	s := NewPackagingService()
	ctx := context.Background()
	invariants := CreateInvariants(250, 500, 1000, 2000, 5000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.NumberOfPacks(ctx, 1000, invariants)
	}
}

func BenchmarkNumberOfPacks_Large(b *testing.B) {
	s := NewPackagingService()
	ctx := context.Background()
	invariants := CreateInvariants(250, 500, 1000, 2000, 5000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.NumberOfPacks(ctx, 50000, invariants)
	}
}

func BenchmarkNumberOfPacks_ManyInvariants(b *testing.B) {
	s := NewPackagingService()
	ctx := context.Background()

	// Create more invariants for a more complex scenario
	packSizes := []int64{10, 25, 50, 100, 250, 500, 750, 1000, 1500, 2000, 3000, 5000}
	invariants := CreateInvariants(packSizes...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.NumberOfPacks(ctx, 10000, invariants)
	}
}
