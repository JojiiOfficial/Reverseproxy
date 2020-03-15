package models

import (
	"net/url"
	"strings"
)

func isURLValid(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u != nil
}

func sliceToLower(s []string) {
	for i := range s {
		s[i] = strings.ToLower(s[i])
	}
}
