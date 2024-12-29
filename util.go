package main

func stringPtr(s string) *string {
	return &s
}

func resolvePtrOrDefault(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
