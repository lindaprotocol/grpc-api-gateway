package repository

import (
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"gorm.io/gorm"
)

// TagRepository struct: Repository for tag operations
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository function: Creates a new tag repository
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// InsertTag function: Inserts a new tag
func (r *TagRepository) InsertTag(tag *models.TagResponse) (int32, error) {
	var id int32
	err := r.db.Raw(`
		INSERT INTO tags (address, tag, description, owner, signature, votes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`, tag.Address, tag.Tag, tag.Description, tag.Owner, "", tag.Votes, tag.CreatedAt).Scan(&id).Error
	return id, err
}

// UpdateTag function: Updates an existing tag
func (r *TagRepository) UpdateTag(id int32, tag, description string) error {
	return r.db.Exec(`
		UPDATE tags 
		SET tag = ?, description = ?
		WHERE id = ?
	`, tag, description, id).Error
}

// DeleteTag function: Deletes a tag
func (r *TagRepository) DeleteTag(id int32) error {
	return r.db.Exec("DELETE FROM tags WHERE id = ?", id).Error
}

// GetTagByID function: Retrieves a tag by ID
func (r *TagRepository) GetTagByID(id int32) (*models.TagResponse, error) {
	var tag models.TagResponse
	err := r.db.Raw(`
		SELECT id, address, tag, description, owner, votes, created_at
		FROM tags
		WHERE id = ?
	`, id).Scan(&tag).Error
	return &tag, err
}

// GetTagsByAddress function: Retrieves tags for an address
func (r *TagRepository) GetTagsByAddress(address string, offset, limit int) ([]*models.TagResponse, int64, error) {
	var tags []*models.TagResponse
	var total int64

	query := r.db.Raw(`
		SELECT id, address, tag, description, owner, votes, created_at
		FROM tags
		WHERE address = ?
		ORDER BY votes DESC, created_at DESC
		OFFSET ? LIMIT ?
	`, address, offset, limit).Scan(&tags)

	r.db.Raw("SELECT COUNT(*) FROM tags WHERE address = ?", address).Scan(&total)

	return tags, total, query.Error
}

// GetAllTags function: Retrieves all tags with pagination
func (r *TagRepository) GetAllTags(offset, limit int, sort string) ([]*models.TagResponse, int64, error) {
	var tags []*models.TagResponse
	var total int64

	orderBy := "votes DESC, created_at DESC"
	if sort != "" {
		if sort[0] == '-' {
			orderBy = sort[1:] + " DESC"
		} else {
			orderBy = sort + " ASC"
		}
	}

	query := r.db.Raw(`
		SELECT id, address, tag, description, owner, votes, created_at
		FROM tags
		ORDER BY ?
		OFFSET ? LIMIT ?
	`, orderBy, offset, limit).Scan(&tags)

	r.db.Raw("SELECT COUNT(*) FROM tags").Scan(&total)

	return tags, total, query.Error
}

// GetRecommendedTags function: Retrieves recommended tags for an address
func (r *TagRepository) GetRecommendedTags(address string, limit int) ([]*models.TagResponse, error) {
	var tags []*models.TagResponse
	
	// Get most popular tags for similar addresses or overall
	query := r.db.Raw(`
		SELECT tag, COUNT(*) as popularity
		FROM tags
		WHERE address IN (
			SELECT DISTINCT address
			FROM tags
			WHERE tag IN (
				SELECT tag
				FROM tags
				WHERE address = ?
				LIMIT 10
			)
		)
		GROUP BY tag
		ORDER BY popularity DESC
		LIMIT ?
	`, address, limit).Scan(&tags)

	return tags, query.Error
}

// VoteTag function: Increments vote count for a tag
func (r *TagRepository) VoteTag(id int32) error {
	return r.db.Exec("UPDATE tags SET votes = votes + 1 WHERE id = ?", id).Error
}