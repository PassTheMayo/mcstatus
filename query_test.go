package mcstatus_test

import (
	"testing"

	"github.com/PassTheMayo/mcstatus"
)

func TestBasicQuery(t *testing.T) {
	_, err := mcstatus.BasicQuery("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}
}

func TestFullQuery(t *testing.T) {
	_, err := mcstatus.FullQuery("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}
}
