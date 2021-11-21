package mcstatus_test

import (
	"testing"

	"github.com/PassTheMayo/mcstatus"
)

func TestBedrockStatus(t *testing.T) {
	_, err := mcstatus.StatusBedrock("localhost", 19132)

	if err != nil {
		t.Fatal(err)
	}
}
