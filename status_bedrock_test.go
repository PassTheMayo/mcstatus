package mcstatus_test

import (
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestBedrockStatus(t *testing.T) {
	_, err := mcstatus.StatusBedrock("grandtheft.mcpe.me", 19132)

	if err != nil {
		t.Fatal(err)
	}
}
