package mcstatus_test

import (
	"fmt"
	"testing"

	mcstatus "github.com/mcstatus-io/MCStatus"
)

func TestStatus(t *testing.T) {
	res, err := mcstatus.Status("play.hypixel.net", 25565)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
