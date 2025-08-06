package helper

import (
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	// Ubah ke huruf kecil
	slug := strings.ToLower(name)

	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile(`[\s\-]+`)
	slug = re.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
