package smartbox

import (
	"encoding/json"
	"testing"
)

func TestStringOrIntUnmarshalInt(t *testing.T) {
	var v stringOrInt
	err := json.Unmarshal([]byte(`42`), &v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 42 {
		t.Errorf("got %d, want 42", v)
	}
}

func TestStringOrIntUnmarshalString(t *testing.T) {
	var v stringOrInt
	err := json.Unmarshal([]byte(`"7"`), &v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 7 {
		t.Errorf("got %d, want 7", v)
	}
}

func TestStringOrIntUnmarshalInvalidString(t *testing.T) {
	var v stringOrInt
	err := json.Unmarshal([]byte(`"abc"`), &v)
	if err == nil {
		t.Fatal("expected error for non-numeric string")
	}
}

func TestStringOrIntUnmarshalInvalid(t *testing.T) {
	var v stringOrInt
	err := json.Unmarshal([]byte(`true`), &v)
	if err == nil {
		t.Fatal("expected error for boolean")
	}
}

func TestStringOrIntInStruct(t *testing.T) {
	type wrapper struct {
		Id stringOrInt `json:"id"`
	}

	tests := []struct {
		name string
		json string
		want stringOrInt
	}{
		{"int value", `{"id": 3}`, 3},
		{"string value", `{"id": "5"}`, 5},
		{"zero int", `{"id": 0}`, 0},
		{"zero string", `{"id": "0"}`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrapper
			if err := json.Unmarshal([]byte(tt.json), &w); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if w.Id != tt.want {
				t.Errorf("got %d, want %d", w.Id, tt.want)
			}
		})
	}
}
