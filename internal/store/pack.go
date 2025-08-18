package store

import (
	"context"
	"errors"

	"github.com/kliuchnikovv/packulator/internal/config"
	"github.com/kliuchnikovv/packulator/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store interface {
	GetPacksInvariantsByHash(ctx context.Context, versionHash string) ([]model.Pack, error)
	SavePack(ctx context.Context, pack *model.Pack) error
	SavePacks(ctx context.Context, packs []model.Pack, versionHash string) error
	GetPackByID(ctx context.Context, id string) (*model.Pack, error)
	ListPacks(ctx context.Context) ([]model.Pack, error)
	GetLatestPackConfig(ctx context.Context) (*model.Pack, error)
	DeletePack(ctx context.Context, id string) error
	HealthCheck(ctx context.Context) error
}

type store struct {
	db *gorm.DB
}

func NewStore(dbConfig *config.DatabaseConfig) (Store, error) {
	db, err := gorm.Open(postgres.Open(dbConfig.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Pack{}, &model.PackItem{}); err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}

func (s *store) SavePack(ctx context.Context, pack *model.Pack) error {
	return s.db.WithContext(ctx).Create(pack).Error
}

func (s *store) SavePacks(ctx context.Context, packs []model.Pack, versionHash string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range packs {
			packs[i].VersionHash = versionHash
			if err := tx.Create(&packs[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *store) GetPackByID(ctx context.Context, id string) (*model.Pack, error) {
	var pack model.Pack
	err := s.db.WithContext(ctx).Preload("PackItems").Where("id = ?", id).First(&pack).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pack, nil
}

func (s *store) GetPacksInvariantsByHash(ctx context.Context, versionHash string) ([]model.Pack, error) {
	var packs []model.Pack
	err := s.db.WithContext(ctx).
		Preload("PackItems").
		Where("version_hash = ?", versionHash).
		Find(&packs).Error

	if err != nil {
		return nil, err
	}

	if len(packs) == 0 {
		return nil, ErrNotFound
	}

	return packs, nil
}

func (s *store) ListPacks(ctx context.Context) ([]model.Pack, error) {
	var packs []model.Pack
	err := s.db.WithContext(ctx).Preload("PackItems").Find(&packs).Error
	return packs, err
}

func (s *store) GetLatestPackConfig(ctx context.Context) (*model.Pack, error) {
	var pack model.Pack
	err := s.db.WithContext(ctx).
		Preload("PackItems").
		Order("created_at DESC").
		First(&pack).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pack, nil
}

func (s *store) DeletePack(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Pack{}).Error
}

func (s *store) HealthCheck(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
