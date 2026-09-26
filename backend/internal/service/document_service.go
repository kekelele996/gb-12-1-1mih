package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// DocumentStore is the persistence port for documents.
type DocumentStore interface {
	Create(d *model.Document) error
	CreateTx(tx *gorm.DB, d *model.Document) error
	FindByID(id uint) (*model.Document, error)
	Update(d *model.Document) error
	UpdateTx(tx *gorm.DB, d *model.Document) error
	ListByApplication(applicationID uint) ([]model.Document, error)
}

// DocumentVersionStore is the persistence port for document versions.
type DocumentVersionStore interface {
	CreateTx(tx *gorm.DB, v *model.DocumentVersion) error
	ListByDocument(documentID uint) ([]model.DocumentVersion, error)
	FindByDocumentAndVersion(documentID uint, versionNo int) (*model.DocumentVersion, error)
}

// AnnotationStore is the persistence port for annotations.
type AnnotationStore interface {
	Create(a *model.Annotation) error
	ListByDocument(documentID uint) ([]model.Annotation, error)
}

// DocumentService handles documents, versions and annotations.
type DocumentService struct {
	db      *gorm.DB
	docRepo DocumentStore
	verRepo DocumentVersionStore
	annRepo AnnotationStore
	appRepo ApplicationProjectStore
	logger  *slog.Logger
}

// NewDocumentService creates a DocumentService.
func NewDocumentService(db *gorm.DB, docRepo DocumentStore, verRepo DocumentVersionStore, annRepo AnnotationStore, appRepo ApplicationProjectStore, logger *slog.Logger) *DocumentService {
	return &DocumentService{db: db, docRepo: docRepo, verRepo: verRepo, annRepo: annRepo, appRepo: appRepo, logger: logger}
}

// loadApplicationForRead locates the parent project and enforces read access.
func (s *DocumentService) loadApplicationForRead(applicationID, userID uint, role string) (*model.ApplicationProject, error) {
	a, err := s.appRepo.FindByID(applicationID)
	if err != nil {
		if errors.Is(err, repositoryErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("ApplicationProject[id=%d] not found", applicationID))
		}
		return nil, fmt.Errorf("document load application: %w", err)
	}
	if err := assertCanViewApplication("ApplicationProject", applicationID, a, userID, role); err != nil {
		return nil, err
	}
	return a, nil
}

// loadDocumentForRead locates a document and enforces read access on its project.
func (s *DocumentService) loadDocumentForRead(id, userID uint, role string) (*model.Document, *model.ApplicationProject, error) {
	d, err := s.docRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositoryErrNotFound) {
			return nil, nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Document[id=%d] not found", id))
		}
		return nil, nil, fmt.Errorf("document find: %w", err)
	}
	a, err := s.loadApplicationForRead(d.ApplicationID, userID, role)
	if err != nil {
		return nil, nil, err
	}
	return d, a, nil
}

// assertStudentCanEdit enforces "owner student + project still editable".
func assertStudentCanEdit(a *model.ApplicationProject, userID uint, role string) error {
	if role != constants.RoleStudent {
		return util.NewAppError(403, constants.CodeForbidden, constants.MsgDocStudentOwnerOnly)
	}
	if a.StudentID != userID {
		return util.NewAppError(403, constants.CodeForbidden, constants.MsgDocStudentOwnerOnly)
	}
	if !constants.IsEditableApplicationStatus(a.Status) {
		return util.NewAppError(409, constants.CodeConflict, constants.MsgDocNotEditable)
	}
	return nil
}

// Create creates a document under an application (student owner, before submission).
func (s *DocumentService) Create(applicationID, userID uint, role, docType, title, content string) (*model.Document, error) {
	if !constants.IsValidDocumentType(docType) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("Document[doc_type=%s] create failed: invalid type", docType))
	}
	a, err := s.loadApplicationForRead(applicationID, userID, role)
	if err != nil {
		return nil, err
	}
	if err := assertStudentCanEdit(a, userID, role); err != nil {
		return nil, err
	}
	d := &model.Document{ApplicationID: applicationID, DocType: docType, Title: title, Content: content, CurrentVersion: 1}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.docRepo.CreateTx(tx, d); err != nil {
			return fmt.Errorf("document create: %w", err)
		}
		if err := s.verRepo.CreateTx(tx, &model.DocumentVersion{DocumentID: d.ID, Content: content, VersionNo: 1, ChangeSummary: "初始版本", CreatedBy: userID}); err != nil {
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

// Save saves content as a new version (student owner, before submission).
func (s *DocumentService) Save(id, userID uint, role, content, changeSummary string) (*model.Document, error) {
	d, a, err := s.loadDocumentForRead(id, userID, role)
	if err != nil {
		return nil, err
	}
	if err := assertStudentCanEdit(a, userID, role); err != nil {
		return nil, err
	}
	d.Content = content
	d.CurrentVersion++
	v := &model.DocumentVersion{DocumentID: id, Content: content, VersionNo: d.CurrentVersion, ChangeSummary: changeSummary, CreatedBy: userID}
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

// Get returns a document by id, enforcing access.
func (s *DocumentService) Get(id, userID uint, role string) (*model.Document, error) {
	d, _, err := s.loadDocumentForRead(id, userID, role)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// ListByApplication returns documents of an application, enforcing access.
func (s *DocumentService) ListByApplication(applicationID, userID uint, role string) ([]model.Document, error) {
	if _, err := s.loadApplicationForRead(applicationID, userID, role); err != nil {
		return nil, err
	}
	return s.docRepo.ListByApplication(applicationID)
}

// ListVersions returns versions of a document, enforcing access.
func (s *DocumentService) ListVersions(documentID, userID uint, role string) ([]model.DocumentVersion, error) {
	if _, _, err := s.loadDocumentForRead(documentID, userID, role); err != nil {
		return nil, err
	}
	return s.verRepo.ListByDocument(documentID)
}

// Rollback restores a document to a historical version (student owner, before submission).
func (s *DocumentService) Rollback(documentID, userID uint, role string, versionNo int) (*model.Document, error) {
	v, err := s.verRepo.FindByDocumentAndVersion(documentID, versionNo)
	if err != nil {
		if errors.Is(err, repositoryErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("DocumentVersion[document_id=%d version=%d] not found", documentID, versionNo))
		}
		return nil, fmt.Errorf("document rollback find: %w", err)
	}
	d, a, err := s.loadDocumentForRead(documentID, userID, role)
	if err != nil {
		return nil, err
	}
	if err := assertStudentCanEdit(a, userID, role); err != nil {
		return nil, err
	}
	d.Content = v.Content
	if err := s.docRepo.Update(d); err != nil {
		return nil, fmt.Errorf("document rollback update: %w", err)
	}
	return d, nil
}

// AddAnnotation adds a counselor/admin annotation (responsible staff only).
func (s *DocumentService) AddAnnotation(counselorID, documentID uint, role, content string, startOffset, endOffset int) (*model.Annotation, error) {
	_, a, err := s.loadDocumentForRead(documentID, counselorID, role)
	if err != nil {
		return nil, err
	}
	if !isManagingStaff(a, counselorID, role) {
		return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgNotResponsibleCounselor)
	}
	an := &model.Annotation{
		DocumentID: documentID, CounselorID: counselorID,
		Content: content, StartOffset: startOffset, EndOffset: endOffset,
	}
	if err := s.annRepo.Create(an); err != nil {
		return nil, fmt.Errorf("annotation create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogAnnotationCreateSuccess, documentID), "id", an.ID)
	return an, nil
}

// ListAnnotations returns annotations of a document, enforcing access.
func (s *DocumentService) ListAnnotations(documentID, userID uint, role string) ([]model.Annotation, error) {
	if _, _, err := s.loadDocumentForRead(documentID, userID, role); err != nil {
		return nil, err
	}
	return s.annRepo.ListByDocument(documentID)
}
