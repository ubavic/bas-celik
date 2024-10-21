package card

import (
	"encoding/hex"
	"slices"
)

type Atr []byte

func (atr Atr) String() string {
	return hex.EncodeToString([]byte(atr))
}

func (atr Atr) Is(otherAtr Atr) bool {
	return slices.Equal(atr, otherAtr)
}

func DetectCardDocumentByAtr(atr Atr) []CardDocumentType {
	if atr.Is(GEMALTO_ATR_1) {
		return []CardDocumentType{GemaltoIdDocumentCardType, VehicleDocumentCardType}
	} else if atr.Is(GEMALTO_ATR_2) || atr.Is(GEMALTO_ATR_3) {
		return []CardDocumentType{GemaltoIdDocumentCardType, MedicalDocumentCardType, VehicleDocumentCardType}
	} else if atr.Is(GEMALTO_ATR_4) {
		return []CardDocumentType{GemaltoIdDocumentCardType}
	} else if atr.Is(MEDICAL_ATR_1) || atr.Is(MEDICAL_ATR_2) {
		return []CardDocumentType{MedicalDocumentCardType}
	} else if atr.Is(VEHICLE_ATR_0) || atr.Is(VEHICLE_ATR_2) || atr.Is(VEHICLE_ATR_3) || atr.Is(VEHICLE_ATR_4) {
		return []CardDocumentType{VehicleDocumentCardType}
	} else if atr.Is(APOLLO_ATR) {
		return []CardDocumentType{ApolloIdDocumentCardType}
	} else {
		return []CardDocumentType{UnknownDocumentCardType}
	}
}
