package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// DocumentService handles documents, versions and annotations.
type DocumentService struct {
	db      *gorm.DB
	docRepo *repository.DocumentRepository
	verRepo *repository.DocumentVersionRepository
	annRepo *repository.AnnotationRepository
	appRepo *repository.ApplicationProjectRepository
	logger  *slog.Logger
}

// NewDocumentService creates a DocumentService.
func NewDocumentService(db *gorm.DB, docRepo *repository.DocumentRepository, verRepo *repository.DocumentVersionRepository, annRepo *repository.AnnotationRepository, appRepo *repository.ApplicationProjectRepository, logger *slog.Logger) *DocumentService {
	return &DocumentService{db: db, docRepo: docRepo, verRepo: verRepo, annRepo: annRepo, appRepo: appRepo, logger: logger}
}

// Create creates a document under an application.
func (s *DocumentService) Create(applicationID uint, docType, title, content string) (*model.Document, error) {
	if !constants.IsValidDocumentType(docType) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Document[doc_type=%s] create failed: invalid type", docType))
	}
	d := &model.Document{ApplicationID: applicationID, DocType: docType, Title: title, Content: content, CurrentVersion: 1}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.docRepo.CreateTx(tx, d); err != nil {
			return fmt.Errorf("document create: %w", err)
		}
		if err := s.verRepo.CreateTx(tx, &model.DocumentVersion{DocumentID: d.ID, Content: content, VersionNo: 1, ChangeSummary: "初始版本"}); err != nil {
			return fmt.Errorf("document initial version create: %w", err)
		}
		return nil
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogDocumentSaveFailed, d.ID), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogDocumentSaveSuccess, d.ID, 1), "application_id", applicationID)
	return d, nil
}

// Save saves content as a new version.
func (s *DocumentService) Save(id uint, content, changeSummary string) (*model.Document, error) {
	d, err := s.docRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Document[id=%d] not found", id))
		}
		return nil, fmt.Errorf("document save find: %w", err)
	}
	d.Content = content
	d.CurrentVersion++
	v := &model.DocumentVersion{DocumentID: id, Content: content, VersionNo: d.CurrentVersion, ChangeSummary: changeSummary}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.docRepo.UpdateTx(tx, d); err != nil {
			return fmt.Errorf("document save update: %w", err)
		}
		if err := s.verRepo.CreateTx(tx, v); err != nil {
			return fmt.Errorf("document version create: %w", err)
		}
		return nil
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogDocumentSaveFailed, id), "error", err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogDocumentSaveSuccess, id, d.CurrentVersion), "version_no", d.CurrentVersion)
	return d, nil
}

// Get returns a document by id.
func (s *DocumentService) Get(id uint) (*model.Document, error) {
	d, err := s.docRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Document[id=%d] not found", id))
		}
		return nil, fmt.Errorf("document get: %w", err)
	}
	return d, nil
}

// ListByApplication returns documents of an application.
func (s *DocumentService) ListByApplication(applicationID uint) ([]model.Document, error) {
	return s.docRepo.ListByApplication(applicationID)
}

// ListVersions returns versions of a document.
func (s *DocumentService) ListVersions(documentID uint) ([]model.DocumentVersion, error) {
	return s.verRepo.ListByDocument(documentID)
}

// Rollback restores a document to a historical version.
func (s *DocumentService) Rollback(documentID uint, versionNo int) (*model.Document, error) {
	v, err := s.verRepo.FindByDocumentAndVersion(documentID, versionNo)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("DocumentVersion[document_id=%d version=%d] not found", documentID, versionNo))
		}
		return nil, fmt.Errorf("document rollback find: %w", err)
	}
	d, err := s.docRepo.FindByID(documentID)
	if err != nil {
		return nil, fmt.Errorf("document rollback doc find: %w", err)
	}
	d.Content = v.Content
	if err := s.docRepo.Update(d); err != nil {
		return nil, fmt.Errorf("document rollback update: %w", err)
	}
	return d, nil
}

// AddAnnotation adds a counselor annotation.
func (s *DocumentService) AddAnnotation(counselorID, documentID uint, content string, startOffset, endOffset int) (*model.Annotation, error) {
	a := &model.Annotation{
		DocumentID: documentID, CounselorID: counselorID,
		Content: content, StartOffset: startOffset, EndOffset: endOffset,
	}
	if err := s.annRepo.Create(a); err != nil {
		return nil, fmt.Errorf("annotation create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogAnnotationCreateSuccess, documentID), "id", a.ID)
	return a, nil
}

// ListAnnotations returns annotations of a document.
func (s *DocumentService) ListAnnotations(documentID uint) ([]model.Annotation, error) {
	return s.annRepo.ListByDocument(documentID)
}
