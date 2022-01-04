package mcstatus_test

import (
	"fmt"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestBasicQuery(t *testing.T) {
	_, err := mcstatus.BasicQuery("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}
}

func TestFullQuery(t *testing.T) {
	v, err := mcstatus.FullQuery("play.dogecraft.net", 25565)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", v)
}
