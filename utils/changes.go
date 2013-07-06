package utils

// This file is an intermediary position between the cacherev task that
// renames files and the dist task, that copies them to the production folder.

var changes map[string]string

// SaveChanges stores a set of name revs for later use.
func SaveChanges(m map[string]string) {
	changes = m
}

// LoadChanges retrieve the last stored set of name revs.
func LoadChanges() map[string]string {
	return changes
}
