package mcstatus_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestBasicQuery(t *testing.T) {
	response, err := mcstatus.BasicQuery("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(response)
}

func TestFullQuery(t *testing.T) {
	response, err := mcstatus.FullQuery("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)
}
