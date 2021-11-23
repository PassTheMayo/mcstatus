package mcstatus

import (
	"testing"
)

func TestDescription(t *testing.T) {
	_, err := NewDescription("\u00A75Test\u00A76Test 2\u00A78Test 3\u00A7kTest 4\u00A7n\u00A7mTest 5")

	if err != nil {
		t.Fatal(err)
	}
}
