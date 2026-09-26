package constants

// MaterialStatus enumerates material checklist statuses.
const (
	MaterialPending  = "pending"
	MaterialUploaded = "uploaded"
	MaterialApproved = "approved"
)

// ValidMaterialStatuses returns all accepted statuses.
func ValidMaterialStatuses() []string {
	return []string{MaterialPending, MaterialUploaded, MaterialApproved}
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
