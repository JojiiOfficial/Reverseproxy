package models

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

// RegexStore store compiled regex to improve performance
type RegexStore struct {
	Store map[string]*regexp.Regexp
}

// NewRegexStore create new regex store
func NewRegexStore() *RegexStore {
	store := RegexStore{}
	store.Store = make(map[string]*regexp.Regexp)
	return &store
}

// GetPattern returns a regexp. If pattern not found, compile and store
func (store *RegexStore) GetPattern(pattern string) *regexp.Regexp {
	if regexp, ok := store.Store[pattern]; ok {
		return regexp
	}

	// Compile pattern
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Error(err)
		return nil
	}

	// Store pattern in RegexStore
	store.Store[pattern] = r

	return r
}
