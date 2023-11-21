package card

import (
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
