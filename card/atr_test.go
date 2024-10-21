package card_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ubavic/bas-celik/card"
)

func Test_DetectCardDocumentByAtr(t *testing.T) {
	testCases := []struct {
		atr            card.Atr
		expectedResult []card.CardDocumentType
	}{
		{
			atr:            card.Atr{},
			expectedResult: []card.CardDocumentType{0},
		},
		{
			atr:            card.Atr{0x01, 0x02},
			expectedResult: []card.CardDocumentType{0},
		},
		{
			atr:            card.APOLLO_ATR,
			expectedResult: []card.CardDocumentType{card.ApolloIdDocumentCardType},
		},
		{
			atr:            card.GEMALTO_ATR_1,
			expectedResult: []card.CardDocumentType{card.GemaltoIdDocumentCardType, card.VehicleDocumentCardType},
		},
		{
			atr:            card.GEMALTO_ATR_2,
			expectedResult: []card.CardDocumentType{card.GemaltoIdDocumentCardType, card.MedicalDocumentCardType, card.VehicleDocumentCardType},
		},
		{
			atr:            card.GEMALTO_ATR_4,
			expectedResult: []card.CardDocumentType{card.GemaltoIdDocumentCardType},
		},
		{
			atr:            card.MEDICAL_ATR_1,
			expectedResult: []card.CardDocumentType{card.MedicalDocumentCardType},
		},
		{
			atr:            card.VEHICLE_ATR_0,
			expectedResult: []card.CardDocumentType{card.VehicleDocumentCardType},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.atr), func(t *testing.T) {
			result := card.DetectCardDocumentByAtr(testCase.atr)
			slices.Sort(result)
			slices.Sort(testCase.expectedResult)

			if !slices.Equal(testCase.expectedResult, result) {
				t.Errorf("Expected response to be %v, but it is %v", testCase.expectedResult, result)
			}
		})

	}
}
