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
			assert.Equal(t, tt.expected, result)
		})
	}
}
