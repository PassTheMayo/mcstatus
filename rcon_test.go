package mcstatus_test

import (
	"testing"

	"github.com/PassTheMayo/mcstatus"
)

func TestRCON(t *testing.T) {
	client := mcstatus.NewRCON()

	if err := client.Dial("localhost", 25575); err != nil {
		t.Fatal(err)
	}

	if err := client.Login("abc123"); err != nil {
		t.Fatal(err)
	}

	if err := client.Run("time query daytime"); err != nil {
		t.Fatal(err)
	}

	if err := client.Run("tps"); err != nil {
		t.Fatal(err)
	}

	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
}
