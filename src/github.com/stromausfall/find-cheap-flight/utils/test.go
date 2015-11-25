package utils

import (
	"fmt"
	"strings"
	"testing"
)

func FailIfStringDoesntHaveSubstring(t *testing.T, s, subString string) {
	if strings.Contains(s, subString) {
		return
	}

	fmt.Println("String : '" + s + "' didn't contain substring : '" + subString + "'")
	t.Fail()
}
