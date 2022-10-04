package internal

import (
	"bytes"
	"encoding/base64"
	"net/url"
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

func URLEncode(input string) string {
	return url.QueryEscape(input)
}

func Base64UrlSafeEncode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}
