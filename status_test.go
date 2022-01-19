package mcstatus_test

import (
	"log"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestStatus(t *testing.T) {
	v, err := mcstatus.Status("play.mc-complex.com", 25565)

	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v\n", v.ModInfo)
}
