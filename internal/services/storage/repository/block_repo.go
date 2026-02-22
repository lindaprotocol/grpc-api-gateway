package repository

import (
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"gorm.io/gorm"
)

type BlockRepository struct {
	db *gorm.DB
}

func NewBlockRepository(db *gorm.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

// SaveBlock saves a block
func (r *BlockRepository) SaveBlock(block *models.Block) error {
	return r.db.Save(block).Error
}

// GetByNumber retrieves a block by number
func (r *BlockRepository) GetByNumber(number int64) (*models.Block, error) {
	var block models.Block
	err := r.db.Where("number = ?", number).First(&block).Error
	return &block, err
}

// GetByHash retrieves a block by hash
func (r *BlockRepository) GetByHash(hash string) (*models.Block, error) {
	var block models.Block
	err := r.db.Where("hash = ?", hash).First(&block).Error
	return &block, err
}

// GetBlocks retrieves paginated blocks
func (r *BlockRepository) GetBlocks(startBlock int64, offset, limit int, sort string) ([]*models.Block, int64, error) {
	var blocks []*models.Block
	var total int64

	query := r.db.Model(&models.Block{})

	if startBlock > 0 {
		query = query.Where("number >= ?", startBlock)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sort != "" {
		if sort[0] == '-' {
			query = query.Order(sort[1:] + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("number DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&blocks).Error; err != nil {
		return nil, 0, err
	}

	return blocks, total, nil
}

// GetLastIndexedBlock returns the last indexed block number
func (r *BlockRepository) GetLastIndexedBlock() (int64, error) {
	var block models.Block
	err := r.db.Order("number DESC").First(&block).Error
	if err != nil {
		return 0, err
	}
	return block.Number, nil
}

// GetBlockRange retrieves a range of blocks
func (r *BlockRepository) GetBlockRange(start, end int64) ([]*models.Block, error) {
	var blocks []*models.Block
	err := r.db.Where("number >= ? AND number <= ?", start, end).
		Order("number ASC").
		Find(&blocks).Error
	return blocks, err
}