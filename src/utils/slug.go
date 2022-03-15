package utils

import "strings"

//SlugOfName is used to enter the slug
func SlugOfName(name string) string {

	name = strings.TrimSpace(name)
	slug := ""
	for i := 0; i < len(name); i++ {
		if !((name[i] > 47 && name[i] < 58) || (name[i] > 64 && name[i] < 91) || (name[i] > 96 && name[i] < 123) || name[i] == 32 || name[i] == 45) {
			continue
		}
		chr := name[i]
		if chr == ' ' {
			chr = '-'
		}
		slug += string(chr)
	}

	return slug
}
