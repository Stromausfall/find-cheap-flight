package utils

import (
	"strings"
	"testing"
	"fmt"
)

func FailIfStringDoesntHaveSubstring(t *testing.T, s, subString string) {
	if strings.Contains(s, subString) {
		return
	}

	fmt.Println("String : '" + s + "' didn't contain substring : '" + subString + "'")
	t.Fail()
}