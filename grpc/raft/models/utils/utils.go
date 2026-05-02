package utils

// GetStringFromPtr safely dereferences a string pointer, returning an empty string if the pointer is nil.
func GetStringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// GetPtrFromString returns a pointer to the given string.
func GetPtrFromString(s string) *string {
	return &s
}
