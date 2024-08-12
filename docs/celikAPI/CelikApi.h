// CelikApi.h
//
#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#ifndef WINAPI
#define WINAPI __stdcall
#endif

#ifndef EID_API
#define EID_API
#endif

#pragma pack(push, 4)

//
// Constants
//

// Size of all UTF-8 and binary fields in bytes

const int EID_MAX_DocRegNo = 9;
const int EID_MAX_DocumentType = 2;
const int EID_MAX_IssuingDate = 10;
const int EID_MAX_ExpiryDate = 10;
const int EID_MAX_IssuingAuthority = 100;
const int EID_MAX_DocumentSerialNumber = 10;
const int EID_MAX_ChipSerialNumber = 14;
const int EID_MAX_DocumentName = 100;

const int EID_MAX_PersonalNumber = 13;
const int EID_MAX_Surname = 200;
const int EID_MAX_GivenName = 200;
const int EID_MAX_ParentGivenName = 200;
const int EID_MAX_Sex = 2;
const int EID_MAX_PlaceOfBirth = 200;
const int EID_MAX_StateOfBirth = 200;
const int EID_MAX_DateOfBirth = 10;
const int EID_MAX_CommunityOfBirth = 200;
const int EID_MAX_StatusOfForeigner = 200;
const int EID_MAX_NationalityFull = 200;
const int EID_MAX_PurposeOfStay = 200;
const int EID_MAX_ENote = 200;

const int EID_MAX_State = 100;
const int EID_MAX_Community = 200;
const int EID_MAX_Place = 200;
const int EID_MAX_Street = 200;
const int EID_MAX_HouseNumber = 20;
const int EID_MAX_HouseLetter = 8;
const int EID_MAX_Entrance = 10;
const int EID_MAX_Floor = 6;
const int EID_MAX_ApartmentNumber = 12;
const int EID_MAX_AddressDate = 10;
const int EID_MAX_AddressLabel = 60;

const int EID_MAX_Portrait = 7700;

const int EID_MAX_Certificate = 2048;

//
// Card types, used in function EidBeginRead
//

const int EID_CARD_ID2008            = 1;
const int EID_CARD_ID2014            = 2;
const int EID_CARD_IF2020            = 3; // ID for foreigners
const int EID_CARD_RP2024            = 4; // Residence permit

//
// Option identifiers, used in function EidSetOption
//

const int EID_O_KEEP_CARD_CLOSED     = 1;

//
// Certificate types, used in function EidReadCertificate
//

const int EID_Cert_MoiIntermediateCA = 1;
const int EID_Cert_User1             = 2;
const int EID_Cert_User2             = 3;

//
// Block types, used in function EidVerifySignature
//

const int EID_SIG_CARD               = 1;
const int EID_SIG_FIXED              = 2;
const int EID_SIG_VARIABLE           = 3;
const int EID_SIG_PORTRAIT           = 4;

// For new card version EidVerifySignature function will return EID_E_UNABLE_TO_EXECUTE for
// parameter EID_SIG_PORTRAIT. Portrait is in new cards part of EID_SIG_FIXED. To determine
// the card version use second parameter of function EidBeginRead

//
// Function return values
//

const int EID_OK                            =  0;
const int EID_E_GENERAL_ERROR               = -1;
const int EID_E_INVALID_PARAMETER           = -2;
const int EID_E_VERSION_NOT_SUPPORTED       = -3;
const int EID_E_NOT_INITIALIZED             = -4;
const int EID_E_UNABLE_TO_EXECUTE           = -5;
const int EID_E_READER_ERROR                = -6;
const int EID_E_CARD_MISSING                = -7;
const int EID_E_CARD_UNKNOWN                = -8;
const int EID_E_CARD_MISMATCH               = -9;
const int EID_E_UNABLE_TO_OPEN_SESSION      = -10;
const int EID_E_DATA_MISSING                = -11;
const int EID_E_CARD_SECFORMAT_CHECK_ERROR  = -12;
const int EID_E_SECFORMAT_CHECK_CERT_ERROR  = -13;
const int EID_E_INVALID_PASSWORD            = -14;
const int EID_E_PIN_BLOCKED                 = -15;

//
// Structures
//

// NOTE: char arrays DO NOT have zero char at the end

typedef struct tagEID_DOCUMENT_DATA
{
	char docRegNo[EID_MAX_DocRegNo];
	int docRegNoSize;
	char documentType[EID_MAX_DocumentType];
	int documentTypeSize;
	char issuingDate[EID_MAX_IssuingDate];
	int issuingDateSize;
	char expiryDate[EID_MAX_ExpiryDate];
	int expiryDateSize;
	char issuingAuthority[EID_MAX_IssuingAuthority];
	int issuingAuthoritySize;
	char documentSerialNumber[EID_MAX_DocumentSerialNumber];
	int documentSerialNumberSize;
	char chipSerialNumber[EID_MAX_ChipSerialNumber];
	int chipSerialNumberSize;
	char documentName[EID_MAX_DocumentName];
	int documentNameSize;
} EID_DOCUMENT_DATA, *PEID_DOCUMENT_DATA;

typedef struct tagEID_FIXED_PERSONAL_DATA
{
	char personalNumber[EID_MAX_PersonalNumber];
	int personalNumberSize;
	char surname[EID_MAX_Surname];
	int surnameSize;
	char givenName[EID_MAX_GivenName];
	int givenNameSize;
	char parentGivenName[EID_MAX_ParentGivenName];
	int parentGivenNameSize;
	char sex[EID_MAX_Sex];
	int sexSize;
	char placeOfBirth[EID_MAX_PlaceOfBirth];
	int placeOfBirthSize;
	char stateOfBirth[EID_MAX_StateOfBirth];
	int stateOfBirthSize;
	char dateOfBirth[EID_MAX_DateOfBirth];
	int dateOfBirthSize;
	char communityOfBirth[EID_MAX_CommunityOfBirth];
	int communityOfBirthSize;
	char statusOfForeigner[EID_MAX_StatusOfForeigner];
	int statusOfForeignerSize;
	char nationalityFull[EID_MAX_NationalityFull];
	int nationalityFullSize;
	char purposeOfStay[EID_MAX_PurposeOfStay];
	int purposeOfStaySize;
	char eNote[EID_MAX_ENote];
	int eNoteSize;
} EID_FIXED_PERSONAL_DATA, *PEID_FIXED_PERSONAL_DATA;

typedef struct tagEID_VARIABLE_PERSONAL_DATA
{
	char state[EID_MAX_State];
	int stateSize;
	char community[EID_MAX_Community];
	int communitySize;
	char place[EID_MAX_Place];
	int placeSize;
	char street[EID_MAX_Street];
	int streetSize;
	char houseNumber[EID_MAX_HouseNumber];
	int houseNumberSize;
	char houseLetter[EID_MAX_HouseLetter];
	int houseLetterSize;
	char entrance[EID_MAX_Entrance];
	int entranceSize;
	char floor[EID_MAX_Floor];
	int floorSize;
	char apartmentNumber[EID_MAX_ApartmentNumber];
	int apartmentNumberSize;
	char addressDate[EID_MAX_AddressDate];
	int addressDateSize;
	char addressLabel[EID_MAX_AddressLabel];
	int addressLabelSize;
} EID_VARIABLE_PERSONAL_DATA, *PEID_VARIABLE_PERSONAL_DATA;

typedef struct tagEID_PORTRAIT
{
	BYTE portrait[EID_MAX_Portrait];
	int portraitSize;
} EID_PORTRAIT, *PEID_PORTRAIT;

typedef struct tagEID_CERTIFICATE
{
	BYTE certificate[EID_MAX_Certificate];
	int certificateSize;
} EID_CERTIFICATE, *PEID_CERTIFICATE;

//
// Functions
//

EID_API int WINAPI EidSetOption(int nOptionID, UINT_PTR nOptionValue);

EID_API int WINAPI EidStartup(int nApiVersion);
EID_API int WINAPI EidCleanup();

EID_API int WINAPI EidBeginRead(LPCSTR szReader, int* pnCardType = 0);
EID_API int WINAPI EidEndRead();

EID_API int WINAPI EidReadDocumentData(PEID_DOCUMENT_DATA pData);
EID_API int WINAPI EidReadFixedPersonalData(PEID_FIXED_PERSONAL_DATA pData);
EID_API int WINAPI EidReadVariablePersonalData(PEID_VARIABLE_PERSONAL_DATA pData);
EID_API int WINAPI EidReadPortrait(PEID_PORTRAIT pData);
EID_API int WINAPI EidReadCertificate(PEID_CERTIFICATE pData, int certificateType);

EID_API int WINAPI EidChangePassword(LPCSTR szOldPassword, LPCSTR szNewPassword, int* pnTriesLeft);
EID_API int WINAPI EidVerifySignature(UINT nSignatureID);


#pragma pack(pop)

#ifdef __cplusplus
};
#endif
