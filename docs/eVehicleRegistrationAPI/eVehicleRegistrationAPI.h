// ---------------------------------
//
#ifndef __eVehicleRegistrationAPI_h__
#define __eVehicleRegistrationAPI_h__

// ---------------------------------

#ifndef SD_API
#define SD_API	extern "C" __declspec(dllimport)
#endif

// ---------------------------------

#define	SD_FIELD(name, length)	\
	char	name[length]; \
	long	name ## Size

// ---------------------------------
//	Structure used to read a registration data file, its
//	signature and the certificate of the document signer
//	These three fields shall be used together to  verify
//	the signature and ensure that the card is not a fake

#define SD_REGISDATA_MAX_SIZE	(4096)
#define SD_SIGNATURE_MAX_SIZE	(1024)
#define SD_AUTHORITY_MAX_SIZE	(4096)

typedef struct groupSD_REGISTRATION_DATA
{
	SD_FIELD( registrationData,		SD_REGISDATA_MAX_SIZE );
	SD_FIELD( signatureData,		SD_SIGNATURE_MAX_SIZE );
	SD_FIELD( issuingAuthority,		SD_AUTHORITY_MAX_SIZE );
} SD_REGISTRATION_DATA;


// ---------------------------------
//	Structure used to read a document data

typedef struct groupSD_DOCUMENT_DATA
{
	SD_FIELD( stateIssuing,					50 );	//	9F33h		--
	SD_FIELD( competentAuthority,			50 );	//	9F35h		--
	SD_FIELD( authorityIssuing,				50 );	//	9F36h		--
	SD_FIELD( unambiguousNumber,			30 );	//	9F38h		--
	SD_FIELD( issuingDate,					16 );	//	8Eh			"I"
	SD_FIELD( expiryDate,					16 );	//	8Dh			"H"
	SD_FIELD( serialNumber,					20 );	//	C9h	(ICCSN)	ext
} SD_DOCUMENT_DATA;

// ---------------------------------
//	Structure used to read a vehicle data

typedef struct groupSD_VEHICLE_DATA
{
	SD_FIELD( dateOfFirstRegistration,		 16 );	//	82h			"B"
	SD_FIELD( yearOfProduction,				  5 );	//	C5h			ext
	SD_FIELD( vehicleMake,					100 );	//	A3h/87h		"D.1"
	SD_FIELD( vehicleType,					100 );	//	A3h/88h		"D.2"
	SD_FIELD( commercialDescription,		100 );	//	A3h/89h		"D.3"
	SD_FIELD( vehicleIDNumber,				100 );	//	8Ah			"E"
	SD_FIELD( registrationNumberOfVehicle,	 20 );	//	81h			"A"
	SD_FIELD( maximumNetPower,				 20 );	//	A5h/91h		"P.2"
	SD_FIELD( engineCapacity,				 20 );	//	A5h/90h		"P.1"
	SD_FIELD( typeOfFuel,					100 );	//	A5h/92h		"P.3"
	SD_FIELD( powerWeightRatio,				 20 );	//	93h			"Q"
	SD_FIELD( vehicleMass,					 20 );	//	8Ch			"G"
	SD_FIELD( maximumPermissibleLadenMass,	 20 );	//	A4h/8Bh		"F.1"
	SD_FIELD( typeApprovalNumber,			 50 );	//	8Fh			"K"
	SD_FIELD( numberOfSeats,				 20 );	//	A6h/94h		"S.1"
	SD_FIELD( numberOfStandingPlaces,		 20 );	//	A6h/95h		"S.2"
	SD_FIELD( engineIDNumber,				100 );	//	A5h/9Eh		"P.5"
	SD_FIELD( numberOfAxles,				 20 );	//	99h			"L"
	SD_FIELD( vehicleCategory,				 50 );	//	98h			"J"
	SD_FIELD( colourOfVehicle,				 50 );	//	9F24h		"R"
	SD_FIELD( restrictionToChangeOwner,		200 );	//	C1h			ext
	SD_FIELD( vehicleLoad,					 20 );	//	C4h			ext
} SD_VEHICLE_DATA;

// ---------------------------------
//	Structure used to read a personal data

typedef struct groupSD_PERSONAL_DATA
{
	SD_FIELD( ownersPersonalNo,				 20 );	//	C2h			ext
	SD_FIELD( ownersSurnameOrBusinessName,	100 );	//	A1h/A7h/83h	"C.1.1"
	SD_FIELD( ownerName,					100 );	//	A1h/A7h/84h	"C.1.2"
	SD_FIELD( ownerAddress,					200 );	//	A1h/A7h/85h	"C.1.3"
	SD_FIELD( usersPersonalNo,				 20 );	//	C3h			ext
	SD_FIELD( usersSurnameOrBusinessName,	100 );	//	A1h/A9h/83h	"C.3.1"
	SD_FIELD( usersName,					100 );	//	A1h/A9h/84h	"C.3.2"
	SD_FIELD( usersAddress,					200 );	//	A1h/A9h/85h	"C.3.3"
} SD_PERSONAL_DATA;

// ---------------------------------
//	Initialization, finalization of library

SD_API long	sdStartup(int apiVersion);
SD_API long	sdCleanup();

// ---------------------------------
//	enumeration, selection of card readers

SD_API long GetReaderName(long index, char* readerName, long* nameSize);
SD_API long SelectReader (char* readerName);

// ---------------------------------
//	new card process request - shall be called prior to read a new card

SD_API long sdProcessNewCard();

// ---------------------------------
//	Read data file, signature & cert for the file with given index
//	(indexes are in 1 up to 3, the 4th file is not signed)

SD_API long sdReadRegistration	(SD_REGISTRATION_DATA*, long index);

// ---------------------------------
//	Read Document, Vehicle and Personal Data

SD_API long sdReadDocumentData	(SD_DOCUMENT_DATA*);
SD_API long sdReadVehicleData	(SD_VEHICLE_DATA*);
SD_API long sdReadPersonalData	(SD_PERSONAL_DATA*);

// ---------------------------------

#endif
