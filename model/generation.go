package model

import (
	"context"
	entity "microservices/entity/model"

	"gorm.io/gorm"
)

type Generation interface {
	Create(ctx context.Context, generation *entity.Generation) error
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	GetByID(ctx context.Context, id uint64) (*entity.Generation, error)
	ListByUserID(ctx context.Context, uid int, page, size int) ([]*entity.Generation, int64, error)
	CreateGenerationFile(ctx context.Context, gf *entity.GenerationFile) error
	GetOutputFileByGenID(ctx context.Context, genID uint64) (*entity.File, error)
	GetFilesByGenerationIDs(ctx context.Context, genIDs []uint64) (map[uint64]map[int]*entity.File, error)
}

type generation struct {
	db *gorm.DB
}

func newGeneration(f *factory) Generation {
	return &generation{
		db: f.db,
	}
}

func (g *generation) Create(ctx context.Context, generation *entity.Generation) error {
	return g.db.WithContext(ctx).Create(generation).Error
}

func (g *generation) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return g.db.WithContext(ctx).Model(&entity.Generation{}).Where("id = ?", id).Updates(updates).Error
}

func (g *generation) GetByID(ctx context.Context, id uint64) (*entity.Generation, error) {
	var gen entity.Generation
	err := g.db.WithContext(ctx).First(&gen, id).Error
	if err != nil {
		return nil, err
	}
	return &gen, nil
}

func (g *generation) ListByUserID(ctx context.Context, uid int, page, size int) ([]*entity.Generation, int64, error) {
	var generations []*entity.Generation
	var total int64

	db := g.db.WithContext(ctx).Model(&entity.Generation{}).Where("user_id = ?", uid)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Order("id desc").Offset(offset).Limit(size).Find(&generations).Error; err != nil {
		return nil, 0, err
	}

	return generations, total, nil
}

func (g *generation) CreateGenerationFile(ctx context.Context, gf *entity.GenerationFile) error {
	return g.db.WithContext(ctx).Create(gf).Error
}

func (g *generation) GetOutputFileByGenID(ctx context.Context, genID uint64) (*entity.File, error) {
	var file entity.File
	// Join generation_files and files tables
	// generation_files has file_id and type. We need type=2 (Output)
	err := g.db.WithContext(ctx).Table("files").
		Joins("JOIN generation_files ON files.id = generation_files.file_id").
		Where("generation_files.generation_id = ? AND generation_files.type = ?", genID, entity.GenerationFileTypeOutput).
		First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (g *generation) GetFilesByGenerationIDs(ctx context.Context, genIDs []uint64) (map[uint64]map[int]*entity.File, error) {
	if len(genIDs) == 0 {
		return nil, nil
	}

	type Result struct {
		GenerationID uint64 `gorm:"column:generation_id"`
		Type         int    `gorm:"column:type"`
		entity.File
	}

	var results []Result
	err := g.db.WithContext(ctx).Table("files").
		Select("files.*, generation_files.generation_id, generation_files.type").
		Joins("JOIN generation_files ON files.id = generation_files.file_id").
		Where("generation_files.generation_id IN ?", genIDs).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Let's return map[uint64]map[int]*entity.File
	// inner map key: 1=Input, 2=Output
	out := make(map[uint64]map[int]*entity.File)
	for _, r := range results {
		if _, ok := out[r.GenerationID]; !ok {
			out[r.GenerationID] = make(map[int]*entity.File)
		}
		f := r.File
		out[r.GenerationID][r.Type] = &f
	}
	return out, nil
}
