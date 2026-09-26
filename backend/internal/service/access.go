package service

import (
	"fmt"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// repositoryErrNotFound is re-exported for services built against store
// interfaces so they do not import gorm repositories directly.
var repositoryErrNotFound = repository.ErrNotFound

// canViewApplication reports whether a user may read a project and its
// sub-resources (documents, materials, timeline).
func canViewApplication(a *model.ApplicationProject, userID uint, role string) bool {
	switch role {
	case constants.RoleStudent:
		return a.StudentID == userID
	case constants.RoleCounselor:
		return true
	default: // admin
		return true
	}
}

// assertCanViewApplication returns a 403/404-style error when access is denied.
func assertCanViewApplication(resType string, resID uint, a *model.ApplicationProject, userID uint, role string) error {
	if canViewApplication(a, userID, role) {
		return nil
	}
	return util.NewAppError(403, constants.CodeForbidden,
		fmt.Sprintf("%s[id=%d] access denied: user_id=%d role=%s is not the student owner", resType, resID, userID, role))
}

// isManagingStaff reports whether a counselor/admin manages this project.
// Admin always manages every project; a counselor must be the responsible
// counselor, unless no counselor has been assigned yet.
func isManagingStaff(a *model.ApplicationProject, userID uint, role string) bool {
	if role != constants.RoleCounselor && role != constants.RoleAdmin {
		return false
	}
	if role == constants.RoleAdmin {
		return true
	}
	return a.CounselorID == 0 || a.CounselorID == userID
}
