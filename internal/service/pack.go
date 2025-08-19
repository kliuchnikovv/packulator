package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/kliuchnikovv/packulator/internal/model"
	"github.com/kliuchnikovv/packulator/internal/store"
)

//go:generate mockgen -source=pack.go -destination=mocks/pack.go -typed

// PackService defines the interface for pack configuration management operations.
type PackService interface {
	// CreatePacks creates a new pack configuration with the given pack sizes
	CreatePacks(ctx context.Context, packs ...int64) (string, error)
	// GetPackByID retrieves a pack configuration by its unique ID
	GetPackByID(ctx context.Context, id string) (*model.Pack, error)
	// GetPackByHash retrieves a pack configuration by its version hash
	GetPackByHash(ctx context.Context, hash string) (*model.Pack, error)
	// ListPacks returns all available pack configurations
	ListPacks(ctx context.Context) ([]model.Pack, error)
	// DeletePack removes a pack configuration by its unique ID
	DeletePack(ctx context.Context, id string) error
}

// packService implements the PackService interface.
type packService struct {
	store store.Store // Database store for pack persistence
}

// NewPackService creates a new pack service instance with the given store.
func NewPackService(store store.Store) PackService {
	return &packService{
		store: store,
	}
}

// CreatePacks creates a new pack configuration from the provided pack sizes.
// It generates a unique version hash and stores the pack configuration in the database.
func (s *packService) CreatePacks(ctx context.Context, packs ...int64) (string, error) {
	// Create pack model with unique ID and version hash
	var pack = model.Pack{
		ID:          uuid.NewString(),
		VersionHash: generateVersionHash(packs),
		PackItems:   make([]model.PackItem, len(packs)),
	}

	// Create pack items for each provided size
	for i, size := range packs {
		pack.TotalAmount += size
		pack.PackItems[i] = model.PackItem{
			ID:     uuid.NewString(),
			PackID: pack.ID,
			Size:   size,
		}
	}

	// Persist pack configuration to database
	err := s.store.SavePacks(ctx, pack)
	if err != nil {
		return "", err
	}

	return pack.VersionHash, nil
}

// GetPackByID retrieves a pack configuration by its unique identifier.
func (s *packService) GetPackByID(ctx context.Context, id string) (*model.Pack, error) {
	return s.store.GetPackByID(ctx, id)
}

// GetPackByHash retrieves a pack configuration by its version hash.
func (s *packService) GetPackByHash(ctx context.Context, hash string) (*model.Pack, error) {
	return s.store.GetPackByHash(ctx, hash)
}

// ListPacks returns all available pack configurations.
func (s *packService) ListPacks(ctx context.Context) ([]model.Pack, error) {
	return s.store.ListPacks(ctx)
}

// DeletePack removes a pack configuration by its unique identifier.
func (s *packService) DeletePack(ctx context.Context, id string) error {
	return s.store.DeletePack(ctx, id)
}

// generateVersionHash creates a deterministic hash from pack sizes.
// It sorts the packs first to ensure the same combination always produces the same hash.
func generateVersionHash(packs []int64) string {
	// Sort packs to ensure deterministic hashing
	sorted := make([]int64, len(packs))
	copy(sorted, packs)
	slices.Sort(sorted)

	// Generate SHA-256 hash from sorted pack sizes
	hash := sha256.New()
	for _, pack := range sorted {
		hash.Write(fmt.Appendf(nil, "%d,", pack))
	}

	// Return first 16 characters of hex-encoded hash
	return fmt.Sprintf("%x", hash.Sum(nil))[:16]
}
