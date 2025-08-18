package model

import (
	"time"

	"gorm.io/gorm"
)

type Pack struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	VersionHash string         `json:"version_hash" gorm:"index;not null"`
	TotalAmount int64          `json:"total_amount" gorm:"not null"`
	PackItems   []PackItem     `json:"pack_items" gorm:"foreignKey:PackID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type PackItem struct {
	ID     string `json:"id" gorm:"primaryKey"`
	PackID string `json:"pack_id" gorm:"not null;index"`
	Size   int64  `json:"size" gorm:"not null"`
}

func (p *Pack) GetPacks() []int64 {
	packs := make([]int64, len(p.PackItems))
	for i, item := range p.PackItems {
		packs[i] = item.Size
	}
	return packs
}
