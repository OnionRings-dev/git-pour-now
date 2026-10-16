// Package semver validates and compares Semantic Versioning 2.0.0 versions.
package semver

import (
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-([0-9A-Za-z-]+\.)*[0-9A-Za-z-]+)?(\+([0-9A-Za-z-]+\.)*[0-9A-Za-z-]+)?$`)

func Valid(value string) bool {
	if !pattern.MatchString(value) {
		return false
	}
	_, pre := parts(value)
	for _, id := range pre {
		if numeric(id) && len(id) > 1 && id[0] == '0' {
			return false
		}
	}
	return true
}

func parts(value string) ([]string, []string) {
	value = strings.SplitN(value, "+", 2)[0]
	sections := strings.SplitN(value, "-", 2)
	core := strings.Split(sections[0], ".")
	if len(sections) == 1 {
		return core, nil
	}
	return core, strings.Split(sections[1], ".")
}

func numeric(value string) bool {
	for _, c := range value {
		if c < '0' || c > '9' {
			return false
		}
	}
	return value != ""
}

func compareNumber(a, b string) int {
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return strings.Compare(a, b)
}

// Compare returns -1, 0 or 1 according to SemVer precedence. Both inputs must
// be valid versions. Numeric identifiers are compared without integer overflow.
func Compare(a, b string) int {
	ac, ap := parts(a)
	bc, bp := parts(b)
	for i := range ac {
		if n := compareNumber(ac[i], bc[i]); n != 0 {
			return n
		}
	}
	if len(ap) == 0 && len(bp) == 0 {
		return 0
	}
	if len(ap) == 0 {
		return 1
	}
	if len(bp) == 0 {
		return -1
	}
	for i := 0; i < len(ap) && i < len(bp); i++ {
		var n int
		an, bn := numeric(ap[i]), numeric(bp[i])
		switch {
		case an && bn:
			n = compareNumber(ap[i], bp[i])
		case an:
			n = -1
		case bn:
			n = 1
		default:
			n = strings.Compare(ap[i], bp[i])
		}
		if n != 0 {
			return n
		}
	}
	if len(ap) < len(bp) {
		return -1
	}
	if len(ap) > len(bp) {
		return 1
	}
	return 0
}
