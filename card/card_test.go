package card

import (
	"bytes"
	"testing"
)

func Test_responseOK(t *testing.T) {
	byteSliceTest := []struct {
		value  []byte
		result bool
	}{{[]byte{0x0F, 0x0F}, false}, {[]byte{0x90, 0x00}, true}}

	for _, testRes := range byteSliceTest {
		res := responseOK(testRes.value)

		if res != testRes.result {
			t.Errorf("Expected res to be %t and it is %t", res, testRes.result)
		}
	}
}

func Test_min(t *testing.T) {
	minTestCases := []struct {
		firstVal  byte
		secondVal byte
		expected  byte
	}{
		{firstVal: 0, secondVal: 0, expected: 0},
		{firstVal: 0x0F, secondVal: 0x1F, expected: 0x0F},
		{firstVal: 0x0F, secondVal: 0x0F, expected: 0x0F},
		{firstVal: 0x0F, secondVal: 0x0C, expected: 0x0C},
	}

	for _, testVal := range minTestCases {
		res := min(int(testVal.firstVal), int(testVal.secondVal))

		if res != int(testVal.expected) {
			t.Errorf("Got %v and we expect %v", res, testVal.expected)
		}
	}

}

func Test_parseTLV(t *testing.T) {

	data := []byte{0x01, 0x00, 0x05, 0x00, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x09, 0x00, 0x05, 0x00, 0x57, 0x6F, 0x72, 0x6C, 0x64}

	testExpectedResults := []struct {
		keyValue     int
		expectResult []byte
	}{
		{keyValue: 1, expectResult: []byte{72, 101, 108, 108, 111}},
		{keyValue: 9, expectResult: []byte{87, 111, 114, 108, 100}},
	}

	res := parseTLV(data)

	if res == nil {
		t.Error("Result should not be null")
	}

	for _, expectRes := range testExpectedResults {
		val, ok := res[uint(expectRes.keyValue)]
		if !ok {
			t.Error("Could not get value")
		}

		if !bytes.Equal(val, expectRes.expectResult) {
			t.Errorf("Expect first element in slice to be %v and we got %v", expectRes.expectResult, val)
		}
	}
}
