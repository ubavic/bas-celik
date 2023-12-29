package card

import (
	"bytes"
	"testing"
)

func Test_responseOK(t *testing.T) {
	byteSliceTest := []struct {
		value  []byte
		result bool
	}{
		{[]byte{0x0F, 0x0F}, false},
		{[]byte{0x90, 0x00}, true},
		{[]byte{0x01, 0xFF, 0x90, 0x00}, true},
		{[]byte{0x01, 0xFF, 0x00, 0x00}, false},
		{[]byte{0xA1}, false},
	}

	for _, testRes := range byteSliceTest {
		res := responseOK(testRes.value)

		if res != testRes.result {
			t.Errorf("Expected res to be %t and it is %t", res, testRes.result)
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

	res, err := parseTLV(data)

	if err != nil {
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

func Test_parseTLV_emptyData(t *testing.T) {

	data := []byte{}
	_, err := parseTLV(data)

	if err == nil {
		t.Error("Error should be raied here - data length is 0")
	}
}

func Test_assignBoolField(t *testing.T) {
	var fields = make(map[uint][]byte)
	var target bool

	var testValues = []struct {
		name   string
		value  []byte
		target bool
	}{
		{name: "correct valuse 0x31", value: []byte{0x31}, target: true},
		{name: "wrong value 0x01", value: []byte{0x01}, target: false},
		{name: "empty value", value: []byte{}, target: false},
	}
	for _, testVal := range testValues {
		fields[0] = testVal.value
		assignBoolField(fields, 0, &target)
		if target != testVal.target {
			t.Fatalf("%s Failed: Expected target to be %v and got %v", testVal.name, testVal.target, target)
		}
	}
}
