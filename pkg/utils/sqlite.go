package utils

import "strings"

func UniqueConstraintFailed(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
