package mcstatus_test

import (
	"fmt"
	"testing"

	"github.com/PassTheMayo/mcstatus/v2"
)

func TestStatusLegacy(t *testing.T) {
	response, err := mcstatus.StatusLegacy("localhost", 25565)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)
}
