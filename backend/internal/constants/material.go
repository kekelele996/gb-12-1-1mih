package constants

// MaterialStatus enumerates material checklist statuses.
const (
	MaterialPending  = "pending"
	MaterialUploaded = "uploaded"
	MaterialApproved = "approved"
	MaterialRejected = "rejected"
)

// ValidMaterialStatuses returns all accepted statuses.
func ValidMaterialStatuses() []string {
	return []string{MaterialPending, MaterialUploaded, MaterialApproved, MaterialRejected}
}

// IsValidMaterialStatus reports whether a status is known.
func IsValidMaterialStatus(s string) bool {
	for _, v := range ValidMaterialStatuses() {
		if v == s {
			return true
		}
	}
	return false
}

// NextMaterialStatuses returns allowed transitions for a material item.
// pending/uploaded/rejected -> uploaded : student (re)uploads the file
// uploaded                 -> approved/rejected : counselor/admin reviews
// approved/rejected        -> uploaded : counselor/admin asks/permits a replacement
func NextMaterialStatuses(s string) []string {
	switch s {
	case MaterialPending:
		return []string{MaterialUploaded}
	case MaterialUploaded:
		return []string{MaterialApproved, MaterialRejected}
	case MaterialApproved, MaterialRejected:
		return []string{MaterialUploaded}
	default:
		return nil
	}
}
