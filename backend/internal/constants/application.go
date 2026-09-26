package constants

// ApplicationStatus enumerates the application project state machine.
const (
	AppStatusPlanning  = "planning"
	AppStatusPreparing = "preparing"
	AppStatusSubmitted = "submitted"
	AppStatusWaiting   = "waiting"
	AppStatusAdmitted  = "admitted"
	AppStatusRejected  = "rejected"
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

// NextApplicationStatuses returns allowed forward transitions.
func NextApplicationStatuses(s string) []string {
	switch s {
	case AppStatusPlanning:
		return []string{AppStatusPreparing}
	case AppStatusPreparing:
		return []string{AppStatusSubmitted}
	case AppStatusSubmitted:
		return []string{AppStatusWaiting}
	case AppStatusWaiting:
		return []string{AppStatusAdmitted, AppStatusRejected, AppStatusWaitlisted}
	default:
		return nil
	}
}
