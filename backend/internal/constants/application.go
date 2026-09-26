package constants

// ApplicationStatus enumerates the application project state machine.
const (
	AppStatusPlanning   = "planning"
	AppStatusPreparing  = "preparing"
	AppStatusSubmitted  = "submitted"
	AppStatusWaiting    = "waiting"
	AppStatusAdmitted   = "admitted"
	AppStatusRejected   = "rejected"
	AppStatusWaitlisted = "waitlisted"
)

// ValidApplicationStatuses returns all accepted statuses.
func ValidApplicationStatuses() []string {
	return []string{
		AppStatusPlanning, AppStatusPreparing, AppStatusSubmitted,
		AppStatusWaiting, AppStatusAdmitted, AppStatusRejected, AppStatusWaitlisted,
	}
}

// IsValidApplicationStatus reports whether a status is known.
func IsValidApplicationStatus(s string) bool {
	for _, v := range ValidApplicationStatuses() {
		if v == s {
			return true
		}
	}
	return false
}

// NextApplicationStatuses returns allowed transitions along the state machine.
// submitted -> preparing is the counselor/admin "return for revision" path.
func NextApplicationStatuses(s string) []string {
	switch s {
	case AppStatusPlanning:
		return []string{AppStatusPreparing}
	case AppStatusPreparing:
		return []string{AppStatusSubmitted}
	case AppStatusSubmitted:
		return []string{AppStatusWaiting, AppStatusPreparing}
	case AppStatusWaiting:
		return []string{AppStatusAdmitted, AppStatusRejected, AppStatusWaitlisted}
	default:
		return nil
	}
}

// EditableApplicationStatuses are the statuses in which a student may still
// edit documents and upload materials. Once submitted, editing is locked
// until the responsible counselor/admin returns the project to preparing.
func EditableApplicationStatuses() map[string]bool {
	return map[string]bool{
		AppStatusPlanning:  true,
		AppStatusPreparing: true,
	}
}

// IsEditableApplicationStatus reports whether a student can edit a project.
func IsEditableApplicationStatus(s string) bool {
	return EditableApplicationStatuses()[s]
}
