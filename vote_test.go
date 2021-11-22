package mcstatus_test

import (
	"testing"
	"time"

	"github.com/PassTheMayo/mcstatus"
)

func TestVote(t *testing.T) {
	err := mcstatus.SendVote("localhost", 8192, mcstatus.VoteOptions{
		ServiceName: "Test",
		Username:    "PassTheMayo",
		Token:       "abc123",
		UUID:        "",
		Timestamp:   time.Now(),
		Timeout:     time.Second * 5,
	})

	if err != nil {
		t.Fatal(err)
	}
}
