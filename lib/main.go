package lib

import (
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validRepoRegex = `^[A-Za-z0-9][A-Za-z0-9\-_]+[A-Za-z0-9]$`
var invalidRepoRegex = `^.*(\-\-|__|\-_|_\-).*$`

func ValidateRepoName(name string) (valid bool, err error) {
	if matchValid, _ := regexp.MatchString(validRepoRegex, name); matchValid {
		if matchInvalid, _ := regexp.MatchString(invalidRepoRegex, name); !matchInvalid {
			return true, nil
		}
	}
	return false, status.Error(codes.InvalidArgument, "Repo name must only contain characters A-Z, a-z, 0-9, -, and _. Repo name must not start or end with - or _. Repo name must not have two or more consecutive - and/or _.")
}
