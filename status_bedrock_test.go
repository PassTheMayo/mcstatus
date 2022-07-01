package mcstatus_test

import (
	"fmt"
	"testing"

	"github.com/PassTheMayo/mcstatus/v3"
)

func TestBedrockStatus(t *testing.T) {
	response, err := mcstatus.StatusBedrock("localhost", 19132)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)
}
