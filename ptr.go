package main

func stringPtr(s string) *string {
	return &s
}

func resolvePtrOrDefault[V any](s *V) V {
	var def V

	if s == nil {
		return def
	}

	return *s
}
