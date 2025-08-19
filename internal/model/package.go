// Package model contains data structures for the Packulator application.
// These models represent database entities and business objects.
package model

import (
	"time"

	"gorm.io/gorm"
)

// Pack represents a pack configuration with its associated pack sizes.
// Each pack configuration has a unique version hash and contains multiple pack items.
type Pack struct {
	ID          string         `json:"id" gorm:"primaryKey"`               // Unique identifier for the pack
	VersionHash string         `json:"version_hash" gorm:"index;not null"` // Version hash for pack configuration
	TotalAmount int64          `json:"total_amount" gorm:"not null"`       // Total amount that can be packed
	PackItems   []PackItem     `json:"pack_items" gorm:"foreignKey:PackID"` // Associated pack items with sizes
	CreatedAt   time.Time      `json:"created_at"`                        // Timestamp when pack was created
	UpdatedAt   time.Time      `json:"updated_at"`                        // Timestamp when pack was last updated
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                     // Soft delete timestamp
}

// PackItem represents an individual pack size within a pack configuration.
// Multiple pack items belong to a single pack configuration.
type PackItem struct {
	ID     string `json:"id" gorm:"primaryKey"`          // Unique identifier for the pack item
	PackID string `json:"pack_id" gorm:"not null;index"` // Foreign key to the parent pack
	Size   int64  `json:"size" gorm:"not null"`          // Size of this pack item
}

// GetPacks extracts and returns all pack sizes from the pack items.
// This is a convenience method for getting just the sizes for calculation algorithms.
func (p *Pack) GetPacks() []int64 {
	packs := make([]int64, len(p.PackItems))
	for i, item := range p.PackItems {
		packs[i] = item.Size
	}
	return packs
}
