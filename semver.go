package main

import (
	"fmt"
	"strings"

	"github.com/aquasecurity/go-version/pkg/semver"
)

const (
	commentDelimiter   = "#"
	tagSlashDelimiter  = "/"
	tagAtmarkDelimiter = "@"
	versionPrefix      = "v"
)

func newSemver(currentTag string, incrLevel string) (string, error) {
	separatedTagWithComment := strings.Split(currentTag, commentDelimiter)
	tagWithoutComment := strings.TrimSpace(separatedTagWithComment[0])
	separatedTag := strings.Split(tagWithoutComment, tagSlashDelimiter)

	var isTagAtmarkDelimiter bool
	isTagAtmarkDelimiter = false
	if len(separatedTag) == 1 {
		separatedTag = strings.Split(tagWithoutComment, tagAtmarkDelimiter)
		if len(separatedTag) > 1 {
			isTagAtmarkDelimiter = true
		}
	}

	currentSemver := separatedTag[len(separatedTag)-1]
	sv, err := semver.Parse(strings.TrimPrefix(currentSemver, versionPrefix))
	if err != nil {
		return "", err
	}
	var nextSemver string
	switch incrLevel {
	case "":
		return "", fmt.Errorf("next_semver_level is empty")
	case "patch":
		nextSemver = sv.IncPatch().String()
	case "minor":
		nextSemver = sv.IncMinor().String()
	case "major":
		nextSemver = sv.IncMajor().String()
	default:
		return "", fmt.Errorf("%s not supported", incrLevel)
	}

	if isTagAtmarkDelimiter {
		separatedTag[len(separatedTag)-1] = fmt.Sprintf("%s", nextSemver)
		return strings.Join(separatedTag, tagAtmarkDelimiter), nil
	}

	separatedTag[len(separatedTag)-1] = fmt.Sprintf("v%s", nextSemver)
	return strings.Join(separatedTag, tagSlashDelimiter), nil
}
