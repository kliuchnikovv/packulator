package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
