// Package store provides database access layer for the Packulator application.
// It implements data persistence operations using GORM ORM with PostgreSQL.
package store

import (
	"context"
	"errors"

	"github.com/kliuchnikovv/packulator/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=pack.go -destination=mocks/pack.go -typed

// Common store errors
var (
	ErrNotFound = errors.New("not found") // Returned when a requested entity is not found
)

// Store defines the interface for database operations on pack configurations.
type Store interface {
	// SavePack persists a single pack configuration to the database
	SavePack(ctx context.Context, pack *model.Pack) error
	// SavePacks persists multiple pack configurations in a single transaction
	SavePacks(ctx context.Context, packs ...model.Pack) error
	// GetPackByID retrieves a pack configuration by its unique ID
	GetPackByID(ctx context.Context, id string) (*model.Pack, error)
	// GetPackByHash retrieves a pack configuration by its version hash
	GetPackByHash(ctx context.Context, hash string) (*model.Pack, error)
	// ListPacks returns all pack configurations in the database
	ListPacks(ctx context.Context) ([]model.Pack, error)
	// DeletePack removes a pack configuration by its unique ID (soft delete)
	DeletePack(ctx context.Context, id string) error
	// HealthCheck verifies database connectivity
	HealthCheck(ctx context.Context) error
}

// store implements the Store interface using GORM ORM.
type store struct {
	db *gorm.DB // GORM database instance
}

// NewStore creates a new store instance with the given GORM dialector.
// It automatically runs database migrations for Pack and PackItem models.
func NewStore(dialector gorm.Dialector) (Store, error) {
	// Initialize GORM database connection
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run automatic database migrations for pack models
	if err := db.AutoMigrate(&model.Pack{}, &model.PackItem{}); err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}

// SavePack persists a single pack configuration to the database.
func (s *store) SavePack(ctx context.Context, pack *model.Pack) error {
	return s.db.WithContext(ctx).Create(pack).Error
}

// SavePacks persists multiple pack configurations in a single database transaction.
// This ensures atomicity - either all packs are saved or none are.
func (s *store) SavePacks(ctx context.Context, packs ...model.Pack) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Save each pack within the transaction
		for i := range packs {
			if err := tx.Create(&packs[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetPackByID retrieves a pack configuration by its unique ID.
// It includes associated PackItems through preloading.
func (s *store) GetPackByID(ctx context.Context, id string) (*model.Pack, error) {
	var pack model.Pack
	// Query pack with preloaded pack items
	err := s.db.WithContext(ctx).Preload("PackItems").Where("id = ?", id).First(&pack).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pack, nil
}

// GetPackByHash retrieves a pack configuration by its version hash.
// It includes associated PackItems through preloading.
func (s *store) GetPackByHash(ctx context.Context, hash string) (*model.Pack, error) {
	var pack model.Pack
	// Query pack by version hash with preloaded pack items
	err := s.db.WithContext(ctx).Preload("PackItems").Where("version_hash = ?", hash).First(&pack).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pack, nil
}

// ListPacks retrieves all pack configurations from the database.
// It includes associated PackItems for each pack through preloading.
func (s *store) ListPacks(ctx context.Context) ([]model.Pack, error) {
	var packs []model.Pack
	// Query all packs with preloaded pack items
	err := s.db.WithContext(ctx).Preload("PackItems").Find(&packs).Error
	return packs, err
}

// DeletePack performs a soft delete on a pack configuration by its unique ID.
// GORM's soft delete sets the DeletedAt timestamp instead of actually removing the record.
func (s *store) DeletePack(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Pack{}).Error
}

// HealthCheck verifies database connectivity by pinging the database.
// This is used by the health check endpoint to ensure the service can connect to the database.
func (s *store) HealthCheck(ctx context.Context) error {
	// Get underlying SQL database instance
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	// Ping database to verify connectivity
	return sqlDB.PingContext(ctx)
}
