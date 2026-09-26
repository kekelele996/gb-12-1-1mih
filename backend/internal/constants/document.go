package constants

// DocumentType enumerates document kinds.
const (
	DocTypePS    = "ps"
	DocTypeRL    = "rl"
	DocTypeCV    = "cv"
	DocTypeEssay = "essay"
)

// ValidDocumentTypes returns all accepted document types.
func ValidDocumentTypes() []string {
	return []string{DocTypePS, DocTypeRL, DocTypeCV, DocTypeEssay}
}

// IsValidDocumentType reports whether a type is known.
func IsValidDocumentType(s string) bool {
	for _, v := range ValidDocumentTypes() {
		if v == s {
			return true
		}
	}
	return false
}
