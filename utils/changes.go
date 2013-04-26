package utils

// This files is an intermediary position between the cacherev task that
// renames files and the dist task, that copies them to the production folder.
// It can be used by the deploy task sometimes to emmit a mapping file too.

var changes map[string]string

func SaveChanges(m map[string]string) {
	changes = m
}

func LoadChanges() map[string]string {
	return changes
}
