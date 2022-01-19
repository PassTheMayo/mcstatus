package mcstatus_test

import (
	"fmt"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestBasicQuery(t *testing.T) {
	_, err := mcstatus.BasicQuery("play.mineluxmc.com", 25565)

	if err != nil {
		t.Fatal(err)
	}
}

func TestFullQuery(t *testing.T) {
	v, err := mcstatus.FullQuery("20.212.168.234", 19132)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", v)
}
