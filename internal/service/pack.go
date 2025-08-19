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

type PackService interface {
	CreatePacks(ctx context.Context, packs ...int64) (string, error)
	GetPackByID(ctx context.Context, id string) (*model.Pack, error)
	GetPackByHash(ctx context.Context, hash string) (*model.Pack, error)
	ListPacks(ctx context.Context) ([]model.Pack, error)
	DeletePack(ctx context.Context, id string) error
}

type packService struct {
	store store.Store
}

func NewPackService(store store.Store) PackService {
	return &packService{
		store: store,
	}
}

func (s *packService) CreatePacks(ctx context.Context, packs ...int64) (string, error) {
	var pack = model.Pack{
		ID:          uuid.NewString(),
		VersionHash: generateVersionHash(packs),
		PackItems:   make([]model.PackItem, len(packs)),
	}
	for i, size := range packs {
		pack.TotalAmount += size
		pack.PackItems[i] = model.PackItem{
			ID:     uuid.NewString(),
			PackID: pack.ID,
			Size:   size,
		}

	}

	err := s.store.SavePacks(ctx, pack)
	if err != nil {
		return "", err
	}

	return pack.VersionHash, nil
}

func (s *packService) GetPackByID(ctx context.Context, id string) (*model.Pack, error) {
	return s.store.GetPackByID(ctx, id)
}

func (s *packService) GetPackByHash(ctx context.Context, hash string) (*model.Pack, error) {
	return s.store.GetPackByHash(ctx, hash)
}

func (s *packService) ListPacks(ctx context.Context) ([]model.Pack, error) {
	return s.store.ListPacks(ctx)
}

func (s *packService) DeletePack(ctx context.Context, id string) error {
	return s.store.DeletePack(ctx, id)
}

func generateVersionHash(packs []int64) string {
	sorted := make([]int64, len(packs))
	copy(sorted, packs)
	slices.Sort(sorted)

	hash := sha256.New()
	for _, pack := range sorted {
		hash.Write([]byte(fmt.Sprintf("%d,", pack)))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))[:16]
}
