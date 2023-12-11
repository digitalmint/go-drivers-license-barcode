package drivers_license

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	expectedDOBT  = time.Date(1969, time.Month(3), 5, 0, 0, 0, 0, time.UTC)
	expectedExpT  = time.Date(2023, time.Month(7), 12, 0, 0, 0, 0, time.UTC)
	expectedDOBT1 = time.Date(1994, time.Month(11), 05, 0, 0, 0, 0, time.UTC)
	expectedExpT1 = time.Date(2026, time.Month(1), 12, 0, 0, 0, 0, time.UTC)

	barcodeTests = []struct {
		str, expectedSerial, expectedDOB, expectedExp string
		expectedDOBT, expectedExpT                    *time.Time
	}{
		// no whitespace separating "DAQ" and the serial number
		{
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgdfgdfxdfg ,sdfzsdfzsdf,D\nDAQFFGG5566\nDBA20230712\nDBB19690305\nasdsda dsadASD F\nasddsaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "20230712",
			expectedExpT:   &expectedExpT,
		},
		// whitespace separating "DAQ" and the serial number
		{
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgdfgdfxdfg ,sdfzsdfzsdf,D\nDAQ\nFFGG5566A\nDBA\n20230712\nDBB\n19690305\nasdsda dsadASD F\nasddsaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566A",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "20230712",
			expectedExpT:   &expectedExpT,
		},
		//  "DBB", and "DBA" exists in data
		{
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgdfgdfxdfg ,sdfzsdfzsdf,D\nDAQ\nFFGG5566B\nDBA\n20230712\nDBB\n19690305\nasds DBBda dsadASD F\nasdds DBAaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566B",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "20230712",
			expectedExpT:   &expectedExpT,
		},
		// invalid expiry date
		{
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgdfgdfxdfg ,sdfzsdfzsdf,D\nDAQ\nFFGG5566C\nDBA\nabc\nDBB\n19690305\nasdsda dsadASD F\nasddsaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566C",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "",
			expectedExpT:   nil,
		},
		// no whitespace in front of "DAQ"
		{
			str:            "@\n\nANSI 636035080002DL00410267ZI03080021DLDAQ3ff15620ed44b4bd2ec27d5d26078a729e\nDCSHyatt\nDDEN\nDACScotty\nDDFN\nDADLamont\nDDGN\nDCAD\nDCBNONE\nDCDNONE\nDBD11051994\nDBB11051994\nDBA01122026\nDBCM\nDAU074 in\nDAYHAZ\nDAG99725 Linda Crossing\nDAIMathilde ton\nDAJKentucky\nDAK786080000\nDCF20190909306CC3856\nDCGUSA\nDAW200\nDCK\nDDAN\nDDB09172015\nZIZIAORG\nZIB\nZIC\nZID",
			expectedSerial: "3ff15620ed44b4bd2ec27d5d26078a729e",
			expectedDOB:    "19941105",
			expectedDOBT:   &expectedDOBT1,
			expectedExp:    "20260112",
			expectedExpT:   &expectedExpT1,
		},
	}
)

func TestDocumentSerialFromBarcodeData(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		require.Equal(t, tt.expectedSerial, bc.DocumentSerial.String)
	}
}

func TestDOBFromBarcodeData(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		require.Equal(t, tt.expectedDOB, bc.Dob.String)
		require.Equal(t, tt.expectedDOBT, bc.Dob.DateT)
	}
}
func TestExpiryFromBarcodeData(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		require.Equal(t, tt.expectedExp, bc.Expiry.String)
		require.Equal(t, tt.expectedExpT, bc.Expiry.DateT)
	}
}

func TestCompareDate(t *testing.T) {
	var bcdata string
	for _, tt := range barcodeTests {
		bc, err := NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		dob, err := bc.SelectDate(BarcodeDataTypeDOB, tt.expectedDOBT)
		require.Nil(t, err)
		require.Equal(t, tt.expectedDOBT, dob)
		exp, err := bc.SelectDate(BarcodeDataTypeExpiry, tt.expectedExpT)
		require.Nil(t, err)
		require.Equal(t, tt.expectedExpT, exp)
	}

	// should return an error
	bcdata = "invalid barcode data"
	_, err := NewBarcode(bcdata)
	require.NotNil(t, err)
	require.IsType(t, ErrInvalidData{}, err)
	require.True(t, IsPackageError(err))
	if !errors.As(err, &ErrInvalidData{}) {
		t.Fatal("invalid error type using errors.As")
	}

	bcdata = "invalid\nbarcode\ndata\nthat passes inspection"
	bc, err := NewBarcode(bcdata)
	require.Nil(t, err)
	require.NotNil(t, bc)
	require.NotNil(t, bc.Dob.Err)
	require.IsType(t, ErrPrefixExtraction{}, bc.Dob.Err)
	require.NotNil(t, bc.Expiry.Err)
	require.IsType(t, ErrPrefixExtraction{}, bc.Expiry.Err)
	require.NotNil(t, bc.DocumentSerial.Err)
	require.IsType(t, ErrPrefixExtraction{}, bc.DocumentSerial.Err)
	//should return true
	require.True(t, IsPackageError(bc.DocumentSerial.Err))
	//should return false
	require.False(t, IsDateError(bc.DocumentSerial.Err))

	// should return nil error and the original DOB
	dobT, err := bc.SelectDate(BarcodeDataTypeDOB, barcodeTests[0].expectedDOBT)
	require.Nil(t, err)
	require.Equal(t, barcodeTests[0].expectedDOBT, dobT)
	// should return nil error and the original exp
	exp, err := bc.SelectDate(BarcodeDataTypeExpiry, barcodeTests[0].expectedExpT)
	require.Nil(t, err)
	require.Equal(t, barcodeTests[0].expectedExpT, exp)

	bc, err = NewBarcode(barcodeTests[0].str)
	require.Nil(t, err)
	require.NotNil(t, bc)

	// should return a ErrBarcodeDateMismatch error and the DOB in the barcode
	passedDOBT := time.Date(2011, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	dobT, err = bc.SelectDate(BarcodeDataTypeDOB, &passedDOBT)
	require.NotNil(t, err)
	require.IsType(t, ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedDOBT, dobT)

	if !errors.As(err, &ErrBarcodeDateMismatch{}) {
		t.Fatal("invalid error type using errors.As")
	}
	//should return true
	require.True(t, IsPackageError(err))
	require.True(t, IsDateError(err))

	// should return a ErrBarcodeDateMismatch error and the expiration date in the barcode
	passedExpT := time.Date(2011, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	dobT, err = bc.SelectDate(BarcodeDataTypeExpiry, &passedExpT)
	require.NotNil(t, err)
	require.IsType(t, ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedExpT, dobT)
	require.EqualValues(t, "fieldname: \"exp\" : barcode date \"20230712\" does not match passed date \"20111231\" - using barcode date", err.Error())

	if !errors.As(err, &ErrBarcodeDateMismatch{}) {
		t.Fatal("invalid error type using errors.As")
	}
	//should return true
	require.True(t, IsPackageError(err))
	require.True(t, IsDateError(err))

	// should return a ErrBarcodeDateMismatch error and the expiration date in the barcode
	passedExpT = time.Time{}
	dobT, err = bc.SelectDate(BarcodeDataTypeExpiry, &passedExpT)
	require.NotNil(t, err)
	require.IsType(t, ErrBarcodeDateMismatch{}, err)
	require.True(t, IsPackageError(err))
	require.Equal(t, barcodeTests[0].expectedExpT, dobT)

	if !errors.As(err, &ErrBarcodeDateMismatch{}) {
		t.Fatal("invalid error type using errors.As")
	}
	//should return true
	require.True(t, IsPackageError(err))
	require.True(t, IsDateError(err))

	// should return a ErrBarcodeDateMismatch error and the expiration date in the barcode
	dobT, err = bc.SelectDate(BarcodeDataTypeExpiry, nil)
	require.NotNil(t, err)
	require.IsType(t, ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedExpT, dobT)

	if !errors.As(err, &ErrBarcodeDateMismatch{}) {
		t.Fatal("invalid error type using errors.As")
	}
	//should return true
	require.True(t, IsPackageError(err))
	require.True(t, IsDateError(err))

	// should return ErrPrefixExtraction
	df := processDate("someteststring", BarcodeDataPrefix("foo"))
	require.NotNil(t, err)
	require.IsType(t, ErrPrefixExtraction{}, df.Err)
	if !errors.As(df.Err, &ErrPrefixExtraction{}) {
		t.Fatal("invalid error type using errors.As")
	}
	require.Nil(t, df.DateT)
	require.Empty(t, df.String)
	require.True(t, IsPackageError(err))

	// should return ErrInvalidDate
	dateT, err := parseDate("someteststring", "testField")
	require.NotNil(t, err)
	var errInvDate ErrInvalidDate
	require.ErrorAs(t, err, &errInvDate)
	require.EqualValues(t, "fieldname: \"testField\" : invalid date: \"someteststring\"", err.Error())
	if !errors.As(err, &ErrInvalidDate{}) {
		t.Fatal("invalid error type using errors.As")
	}
	require.Nil(t, dateT)
	//should return true
	require.True(t, IsPackageError(err))
	require.True(t, IsDateError(err))

	// should return ErrParseDate
	dateT, err = parseDate("99999999", "testField")
	require.NotNil(t, err)
	require.EqualValues(t, "fieldname: \"testField\" : could not parse date \"99999999\"", err.Error())
	require.IsType(t, ErrParseDate{}, err)
	require.Nil(t, dateT)
	if !errors.As(err, &ErrParseDate{}) {
		t.Fatal("invalid error type using errors.As")
	}

	//should return true
	require.True(t, IsPackageError(err))
	require.True(t, IsDateError(err))
}
