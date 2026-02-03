package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	slug, _, _ = transform.String(t, slug)

	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")

	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

func GenerateUniqueSlug(title string, suffix string) string {
	slug := GenerateSlug(title)
	if suffix != "" {
		slug = slug + "-" + suffix
	}
	return slug
}
