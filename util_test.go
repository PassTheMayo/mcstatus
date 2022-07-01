package mcstatus_test

import (
	"log"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestParseAddress(t *testing.T) {
	host, port, err := mcstatus.ParseAddress("play.hypixel.net", 25565)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(host)
	log.Println(port)
}
