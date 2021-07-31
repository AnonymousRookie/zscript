package utils

import (
	"testing"
)

func TestUtils(t *testing.T) {
	if IsNumeric('a') == true {
		t.Error("'a' isNoCharNumeric!")
	}
	if IsNumeric('1') != true {
		t.Error("'1' isCharNumeric!")
	}
	if IsDelimiter('[') != true {
		t.Error("'[' isCharDelimiter!")
	}
}
