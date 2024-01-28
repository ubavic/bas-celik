package card

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func Test_parseTLV(t *testing.T) {
	testCases := []struct {
		data           []byte
		expectedResult map[uint][]byte
		expectedError  error
	}{
		{
			data:          []byte{},
			expectedError: ErrInvalidLength,
		},
		{
			data: []byte{0x01, 0x00, 0x05, 0x00, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x09, 0x00, 0x05, 0x00, 0x57, 0x6F, 0x72, 0x6C, 0x64},
			expectedResult: map[uint][]byte{
				9: {87, 111, 114, 108, 100},
				1: {72, 101, 108, 108, 111},
			},
			expectedError: nil,
		},
		{
			data:          []byte{0x01, 0x00, 0x05, 0x00, 0x48, 0x65},
			expectedError: ErrInvalidLength,
		},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Case %d", i),
			func(t *testing.T) {
				result, err := parseTLV(testCase.data)

				if err == nil && testCase.expectedError == nil {
					if !reflect.DeepEqual(result, testCase.expectedResult) {
						t.Errorf("Result is not expected. Got\n%v\nbut expected\n%v", result, testCase.expectedResult)
					}
				} else {
					if !errors.Is(err, testCase.expectedError) {
						t.Errorf("Expected error %v, but got %v", testCase.expectedError, err)
					}
				}
			},
		)
	}
}

func Test_assignField(t *testing.T) {
	t.Run(
		"Case 1",
		func(t *testing.T) {
			var target string

			fields := make(map[uint][]byte)
			fields[0] = []byte{}

			assignField[uint](fields, 0, &target)
			if target != "" {
				t.Fatalf("Expected an empty string, but got '%s'", target)
			}
		},
	)

	t.Run(
		"Case 2",
		func(t *testing.T) {
			var target string

			fields := make(map[uint][]byte)
			fields[5] = []byte("Hello")

			assignField[uint](fields, 5, &target)
			if target != "Hello" {
				t.Fatalf("Expected 'Hello', but got '%s'", target)
			}
		},
	)

	t.Run(
		"Case 3",
		func(t *testing.T) {
			var target string

			fields := make(map[uint][]byte)

			assignField[uint](fields, 0, &target)
			if target != "" {
				t.Fatalf("Expected an empty string, but got '%s'", target)
			}
		},
	)
}

func Test_assignBoolField(t *testing.T) {
	var testCases = []struct {
		value  []byte
		target bool
	}{
		{value: []byte{0x31}, target: true},
		{value: []byte{0x01}, target: false},
		{value: []byte{}, target: false},
	}

	for i, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Case %d", i),
			func(t *testing.T) {
				var target bool

				fields := make(map[uint][]byte)
				fields[0] = testCase.value

				assignBoolField(fields, 0, &target)
				if target != testCase.target {
					t.Fatalf("Expected target to be %v, but got %v", testCase.target, target)
				}
			},
		)
	}
}
