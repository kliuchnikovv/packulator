package service

import (
	"context"
	"errors"
	"maps"
	"sort"
)

type variant struct {
	numberOfPacks int64
	overshoot     int64
	combination   map[int64]int64
}

func NumberOfPacks(
	ctx context.Context,
	amount int64,
	packs []int64,
) (map[int64]int64, error) {
	if amount <= 0 || len(packs) == 0 {
		return map[int64]int64{}, nil
	}

	sort.Slice(packs, func(i, j int) bool {
		return packs[i] < packs[j]
	})

	var (
		max      = amount + packs[len(packs)-1]
		variants = make([]*variant, max+1)
	)

	variants[0] = &variant{
		numberOfPacks: 0,
		overshoot:     -amount,
		combination:   make(map[int64]int64),
	}

	// Dynamic programming: build all possible combinations
	for sum := int64(0); sum <= max; sum++ {
		if variants[sum] == nil {
			continue
		}

		// Try adding each available pack size
		for _, pack := range packs {
			var newSum = sum + pack
			if newSum > max {
				break // Optimization: packs are sorted, so we can break early
			}
			current := variants[sum]

			// Create new variant by adding current pack to existing combination
			newVariant := &variant{
				numberOfPacks: current.numberOfPacks + 1,
				overshoot:     newSum - amount,
				combination:   make(map[int64]int64, len(current.combination)),
			}

			maps.Copy(newVariant.combination, current.combination)
			newVariant.combination[pack]++

			// Keep the better variant for this sum
			if variants[newSum] == nil {
				variants[newSum] = newVariant
			} else if isBetter(newVariant, variants[newSum]) {
				variants[newSum] = newVariant
			}
		}
	}

	var result = getOptimalVariant(amount, max, variants)
	if result == nil {
		return nil, errors.New("could not find a valid combination")
	}

	return result.combination, nil
}

func getOptimalVariant(amount, max int64, variants []*variant) *variant {
	var result *variant
	for s := amount; s <= max; s++ {
		if variants[s] == nil {
			continue
		}
		if result == nil || isBetter(variants[s], result) {
			result = variants[s]
		}
	}

	return result
}

func isBetter(left, right *variant) bool {
	if left.overshoot < right.overshoot {
		return true
	} else if left.overshoot > right.overshoot {
		return false
	}

	return left.numberOfPacks < right.numberOfPacks
}
