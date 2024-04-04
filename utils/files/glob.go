package files

import (
	"path/filepath"
	"sort"

	log "github.com/sirupsen/logrus"
)

func GlobFiles(globs []string) []string {
	filesSet := map[string]bool{}

	for _, glob := range globs {
		files, err := filepath.Glob(glob)
		if err != nil {
			log.WithField("glob", glob).Warnf("failed to glob path, skipping...")
			continue
		}

		for _, file := range files {
			filesSet[file] = true
		}
	}

	files := []string{}
	for file := range filesSet {
		files = append(files, file)
	}
	sort.Strings(files)

	return files
}
