package mcstatus_test

import (
	"fmt"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestStatus(t *testing.T) {
	response, err := mcstatus.Status("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)
}
