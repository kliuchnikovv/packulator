package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/kliuchnikovv/packulator/internal/service"
)

func TestDebugInvariants(t *testing.T) {
	invariants := service.CreateInvariants(250, 500, 1000, 2000, 5000)
	
	fmt.Printf("Number of invariants: %d\n", len(invariants))
	for i, inv := range invariants {
		fmt.Printf("Invariant %d: TotalAmount=%d, Items=%d\n", i, inv.TotalAmount, len(inv.PackItems))
		for j, item := range inv.PackItems {
			fmt.Printf("  Item %d: Size=%d\n", j, item.Size)
		}
	}
	
	s := service.NewPackagingService()
	result, err := s.NumberOfPacks(context.Background(), 3750, invariants)
	if err != nil {
		t.Fatal(err)
	}
	
	fmt.Printf("Result for amount 3750: %+v\n", result)
	for packSize, count := range result {
		fmt.Printf("Pack size %d: count %d\n", packSize, count)
	}
}