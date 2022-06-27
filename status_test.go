package mcstatus_test

import (
	"log"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestStatus(t *testing.T) {
	response, err := mcstatus.Status("play.hypixel.net", 25565)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(response)
}
