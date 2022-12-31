package slices

// Contains reports whether v is within s.
func Contains[T comparable](s []T, v T) bool {
	return Index(s, v) >= 0
}

// Index returns the index of the first instance of v in s, or -1 if v is not
// present in s.
func Index[T comparable](s []T, v T) int {
	for i := range s {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// Reverse reorders the elements of s in place so that any two elements of s
// that were previously in a certain order are now in the opposite order.
func Reverse[T any](s []T) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-1-i] = s[len(s)-1-i], s[i]
	}
}
