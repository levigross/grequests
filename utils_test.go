package grequests

import (
	"testing"
)

func TestConvertInterfaceToString(t *testing.T) {
	if ConvertIToString("one") != "one" {
		t.Error("Error converting string")
	}
	if ConvertIToString(`one`) != "one" {
		t.Error("Error converting string")
	}
	if ConvertIToString(int(123)) != "123" {
		t.Error("Error converting int")
	}
	if ConvertIToString(int8(123)) != "123" {
		t.Error("Error converting int8")
	}
	if ConvertIToString(int16(123)) != "123" {
		t.Error("Error converting int16")
	}
	if ConvertIToString(int32(123)) != "123" {
		t.Error("Error converting int32")
	}
	if ConvertIToString(int64(123)) != "123" {
		t.Error("Error converting int64")
	}
	if ConvertIToString(uint(123)) != "123" {
		t.Error("Error converting uint")
	}
	if ConvertIToString(uint8(123)) != "123" {
		t.Error("Error converting uint8")
	}
	if ConvertIToString(uint16(123)) != "123" {
		t.Error("Error converting uint16")
	}
	if ConvertIToString(uint32(123)) != "123" {
		t.Error("Error converting uint32")
	}
	if ConvertIToString(uint64(123)) != "123" {
		t.Error("Error converting uint64")
	}
//	if ConvertIToString('⌘') != "⌘" {
//		t.Error("Error converting rune got: ", ConvertIToString('⌘'))
//	}
}
