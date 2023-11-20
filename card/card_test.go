package card

import (
	"testing"

	"github.com/ebfe/scard"
)

func TestReadCard(t *testing.T) {
	var card scard.Card
	_, err := ReadCard(&card)

	if err == nil {
		t.Errorf("Expected error here!")
	}
}
