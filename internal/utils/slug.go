package utils

import (
	"regexp"
	"strings"
)

func Slugify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "-")

	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}

func GenerateUniqueSlug(title string, existingSlugs []string) string {
	baseSlug := Slugify(title)

	isUnique := true
	for _, existing := range existingSlugs {
		if existing == baseSlug {
			isUnique = false
			break
		}
	}

	if isUnique {
		return baseSlug
	}

	counter := 1
	for {
		newSlug := baseSlug + "-" + string(rune('0'+counter))
		isUnique = true
		for _, existing := range existingSlugs {
			if existing == newSlug {
				isUnique = false
				break
			}
		}
		if isUnique {
			return newSlug
		}
		counter++
	}
}
