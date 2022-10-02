package internal

import (
	"bytes"
	"sort"
)

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func InsertNth(s string, r rune, n int) string {
	var buffer bytes.Buffer

	var n1 = n - 1

	var l1 = len(s) - 1

	for i, rune := range s {
		buffer.WriteRune(rune)

		if i%n == n1 && i != l1 {
			buffer.WriteRune(r)
		}
	}

	return buffer.String()
}
