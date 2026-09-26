package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// DocumentRepository handles document persistence.
type DocumentRepository struct{ db *gorm.DB }

// NewDocumentRepository creates the repository.
func NewDocumentRepository(db *gorm.DB) *DocumentRepository { return &DocumentRepository{db: db} }

// Create inserts a document.
func (r *DocumentRepository) Create(d *model.Document) error { return translate(r.db.Create(d).Error) }

// CreateTx inserts a document within an outer transaction.
func (r *DocumentRepository) CreateTx(tx *gorm.DB, d *model.Document) error {
	return translate(tx.Create(d).Error)
}

// FindByID locates a document by id.
func (r *DocumentRepository) FindByID(id uint) (*model.Document, error) {
	var d model.Document
	if err := translate(r.db.First(&d, id).Error); err != nil {
		return nil, err
	}
	return &d, nil
}

// Update persists a document.
func (r *DocumentRepository) Update(d *model.Document) error { return translate(r.db.Save(d).Error) }

// UpdateTx persists a document within an outer transaction.
func (r *DocumentRepository) UpdateTx(tx *gorm.DB, d *model.Document) error {
	return translate(tx.Save(d).Error)
}

// ListByApplication returns documents of an application.
func (r *DocumentRepository) ListByApplication(applicationID uint) ([]model.Document, error) {
	var items []model.Document
	if err := r.db.Where("application_id = ?", applicationID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
