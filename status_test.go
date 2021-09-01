package mcstatus_test

import (
	"testing"

	mcstatus "github.com/mcstatus-io/MCStatus"
)

func TestStatus(t *testing.T) {
	_, err := mcstatus.Status("play.hypixel.net", 25565)

	if err != nil {
		t.Fatal(err)
	}
}
