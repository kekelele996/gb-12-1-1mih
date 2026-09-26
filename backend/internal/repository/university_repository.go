package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// UniversityRepository handles university persistence.
type UniversityRepository struct{ db *gorm.DB }

// NewUniversityRepository creates the repository.
func NewUniversityRepository(db *gorm.DB) *UniversityRepository { return &UniversityRepository{db: db} }

// Create inserts a university.
func (r *UniversityRepository) Create(u *model.University) error { return translate(r.db.Create(u).Error) }

// FindByID locates a university by id.
func (r *UniversityRepository) FindByID(id uint) (*model.University, error) {
	var u model.University
	if err := translate(r.db.First(&u, id).Error); err != nil {
		return nil, err
	}
	return &u, nil
}

// Update persists a university.
func (r *UniversityRepository) Update(u *model.University) error { return translate(r.db.Save(u).Error) }

// Delete removes a university.
func (r *UniversityRepository) Delete(id uint) error {
	res := r.db.Delete(&model.University{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List filters universities by country/ranking/major with pagination.
func (r *UniversityRepository) List(country string, rankMin, rankMax int, keyword string, page, pageSize int) ([]model.University, int64, error) {
	var items []model.University
	var total int64
	q := r.db.Model(&model.University{})
	if country != "" {
		q = q.Where("country = ?", country)
	}
	if rankMax > 0 {
		q = q.Where("ranking >= ? AND ranking <= ?", rankMin, rankMax)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR top_majors LIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("ranking ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListByIDs returns universities matching the given ids (for recommendations).
func (r *UniversityRepository) ListByIDs(ids []uint) ([]model.University, error) {
	var items []model.University
	if len(ids) == 0 {
		return items, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
