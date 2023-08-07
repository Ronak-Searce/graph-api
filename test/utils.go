package test

import (
	"strings"

	"github.com/google/uuid"
)

var hyphenRemover = strings.NewReplacer("-", "")

func getRandomEmail() string {
	return hyphenRemover.Replace(uuid.New().String())[:25] + "@dummy.com"
}

func getRandomString() string {
	return hyphenRemover.Replace(uuid.New().String())[:25]
}
