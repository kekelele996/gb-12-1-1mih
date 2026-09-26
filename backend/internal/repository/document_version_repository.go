package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// DocumentVersionRepository handles version persistence.
type DocumentVersionRepository struct{ db *gorm.DB }

// NewDocumentVersionRepository creates the repository.
func NewDocumentVersionRepository(db *gorm.DB) *DocumentVersionRepository {
	return &DocumentVersionRepository{db: db}
}

// Create inserts a version.
func (r *DocumentVersionRepository) Create(v *model.DocumentVersion) error {
	return translate(r.db.Create(v).Error)
}

// CreateTx inserts a version within an outer transaction.
func (r *DocumentVersionRepository) CreateTx(tx *gorm.DB, v *model.DocumentVersion) error {
	return translate(tx.Create(v).Error)
}

// ListByDocument returns versions of a document ordered by version number desc.
func (r *DocumentVersionRepository) ListByDocument(documentID uint) ([]model.DocumentVersion, error) {
	var items []model.DocumentVersion
	if err := r.db.Where("document_id = ?", documentID).Order("version_no DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByDocumentAndVersion locates a specific version.
func (r *DocumentVersionRepository) FindByDocumentAndVersion(documentID uint, versionNo int) (*model.DocumentVersion, error) {
	var v model.DocumentVersion
	if err := translate(r.db.Where("document_id = ? AND version_no = ?", documentID, versionNo).First(&v).Error); err != nil {
		return nil, err
	}
	return &v, nil
}
