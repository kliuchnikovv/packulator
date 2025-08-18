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

type PackService interface {
	CreatePacks(ctx context.Context, packs ...int64) (string, error)
	GetPacksByVersionHash(ctx context.Context, versionHash string) ([]model.Pack, error)
	GetPackByID(ctx context.Context, id string) (*model.Pack, error)
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
	versionHash := generateVersionHash(packs)

	invariants := CreateInvariants(packs...)

	modelPacks := make([]model.Pack, len(invariants))
	for i, inv := range invariants {
		modelPacks[i] = model.Pack{
			ID:          uuid.NewString(),
			VersionHash: versionHash,
			TotalAmount: inv.TotalAmount,
			PackItems:   make([]model.PackItem, len(inv.PackItems)),
		}

		for j, pack := range inv.PackItems {
			modelPacks[i].PackItems[j] = model.PackItem{
				ID:     uuid.NewString(),
				PackID: modelPacks[i].ID,
				Size:   pack.Size,
			}
		}
	}

	err := s.store.SavePacks(ctx, modelPacks, versionHash)
	if err != nil {
		return "", err
	}

	return versionHash, nil
}

func (s *packService) GetPacksByVersionHash(ctx context.Context, versionHash string) ([]model.Pack, error) {
	return s.store.GetPacksInvariantsByHash(ctx, versionHash)
}

func (s *packService) GetPackByID(ctx context.Context, id string) (*model.Pack, error) {
	return s.store.GetPackByID(ctx, id)
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
