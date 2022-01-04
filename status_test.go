package mcstatus_test

import (
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestStatus(t *testing.T) {
	_, err := mcstatus.Status("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}
}
