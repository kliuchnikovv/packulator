package service

import (
	"sort"

	"github.com/google/uuid"
	"github.com/kliuchnikovv/packulator/internal/model"
)

func CreateInvariants(packs ...int64) []model.Pack {
	sort.Slice(packs, func(i int, j int) bool {
		return packs[i] < packs[j]
	})

	var result = make([]model.Pack, 0, len(packs))
	var previousInvariants = make([]model.Pack, 0)

	for _, size := range packs {
		var packID = uuid.NewString()
		result = append(result, model.Pack{
			ID: packID,
			PackItems: []model.PackItem{{
				ID:     uuid.NewString(),
				PackID: packID,
				Size:   size,
			}},
			TotalAmount: size,
		})

		for i := len(previousInvariants) - 1; i >= 0; i-- {
			var newInvariant = model.Pack{
				ID: uuid.NewString(),
			}

			for j := 0; j < len(previousInvariants[i].PackItems); j++ {
				newInvariant.PackItems = append(newInvariant.PackItems, previousInvariants[i].PackItems[j])
			}

			newInvariant.PackItems = append(newInvariant.PackItems, model.PackItem{
				ID:     uuid.NewString(),
				PackID: newInvariant.ID,
				Size:   size,
			})
			newInvariant.TotalAmount = previousInvariants[i].TotalAmount + size

			result = append(result, newInvariant)
		}

		previousInvariants = result[:]
	}

	sortInvariants(result, true)

	return result
}

func sortInvariants(invariants []model.Pack, ascending bool) {
	var cmpFunc = func(a, b int64) bool {
		return a > b
	}
	if ascending {
		cmpFunc = func(a, b int64) bool {
			return a < b
		}
	}

	sort.Slice(invariants, func(i, j int) bool {
		if invariants[i].TotalAmount == invariants[j].TotalAmount {
			return len(invariants[i].PackItems) < len(invariants[j].PackItems)
		}

		return cmpFunc(invariants[i].TotalAmount, invariants[j].TotalAmount)
	})
}
