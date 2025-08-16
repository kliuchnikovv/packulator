package service

import (
	"context"
	"math"
	"sort"

	"github.com/kliuchnikovv/packulator/internal/model"
)

type PackagingService interface {
	NumberOfPacks(ctx context.Context, amount int64, invariants []model.Pack) (map[int64]int64, error)
}

type service struct {
}

func NewPackagingService() PackagingService {
	return &service{}
}

func (s *service) NumberOfPacks(
	ctx context.Context,
	amount int64,
	invariants []model.Pack,
) (map[int64]int64, error) {
	sortInvariants(invariants, false)

	result := make(map[int64]int64)
	for i, invariant := range invariants {
		if amount <= 0 {
			break
		}

		if i < len(invariants)-1 && invariants[i+1].TotalAmount >= amount {
			continue
		}

		sort.Slice(invariant.PackItems, func(i, j int) bool {
			return invariant.PackItems[i].Size > invariant.PackItems[j].Size
		})

		for j, item := range invariant.PackItems {
			if amount <= 0 {
				break
			}

			var numberOfPacks = float64(amount) / float64(item.Size)

			if j < len(invariant.PackItems)-1 {
				if numberOfPacks <= float64(invariant.PackItems[j+1].Size)/float64(item.Size) {
					continue
				}
			}

			if numberOfPacks < 1 {
				result[item.Size] = int64(math.Round(numberOfPacks))
			} else {
				result[item.Size] = int64(math.Floor(numberOfPacks))
			}

			amount -= result[item.Size] * item.Size
		}

		var smallestPackSize = invariant.PackItems[len(invariant.PackItems)-1].Size
		if amount > 0 && amount < smallestPackSize {
			result[smallestPackSize]++
			amount -= smallestPackSize
		}

	}

	return result, nil
}
