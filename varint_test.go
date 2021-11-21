package mcstatus_test

import (
	"reflect"
	"testing"

	"github.com/PassTheMayo/mcstatus"
)

type test struct {
	Result int64
	Bytes  []byte
}

var (
	tests = []test{
		{
			Result: 0,
			Bytes:  []byte{0x00},
		},
		{
			Result: 1,
			Bytes:  []byte{0x01},
		},
		{
			Result: 2,
			Bytes:  []byte{0x02},
		},
		{
			Result: 127,
			Bytes:  []byte{0x7F},
		},
		{
			Result: 128,
			Bytes:  []byte{0x80, 0x01},
		},
		{
			Result: 255,
			Bytes:  []byte{0xFF, 0x01},
		},
		{
			Result: 2097151,
			Bytes:  []byte{0xFF, 0xFF, 0x7F},
		},
		{
			Result: 2147483647,
			Bytes:  []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07},
		},
		{
			Result: -1,
			Bytes:  []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F},
		},
		{
			Result: -2147483648,
			Bytes:  []byte{0x80, 0x80, 0x80, 0x80, 0x08},
		},
	}
)

func TestReadVarInt(t *testing.T) {
	for _, v := range tests {
		p := mcstatus.NewPacket()

		if err := p.WriteBytes(v.Bytes); err != nil {
			t.Fatal(err)
		}

		res, err := p.ReadVarInt()

		if err != nil {
			t.Fatal(err)
		}

		if res != int32(v.Result) {
			t.Fatalf("read varint test failed, expected %d, got %d", v.Result, res)
		}
	}

	for _, v := range tests {
		p := mcstatus.NewPacket()

		if err := p.WriteVarInt(int32(v.Result)); err != nil {
			t.Fatal(err)
		}

		data, err := p.ReadAllBytes()

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(v.Bytes, data) {
			t.Fatalf("write varint test failed, expected %+v, got %+v", v.Bytes, data)
		}
	}
}
