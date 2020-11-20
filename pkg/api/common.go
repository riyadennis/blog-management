package api

import (
	"path"
	"strings"
)

// ParsePath splits off domain name from the URL
// head wil have the root component and
// tail will have rest of the url path
// p = localhost:8080/api/v1/article/{REF-ID} => head: api & tail v1/article/{REF-ID}
func ParsePath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
