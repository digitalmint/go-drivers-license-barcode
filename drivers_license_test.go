package drivers_license_test

import (
	"errors"
	"testing"
	"time"

	drivers_license "github.com/digitalmint/go-drivers-license-barcode"
	"github.com/stretchr/testify/require"
)

var (
	expectedDOBT = time.Date(1969, time.Month(3), 5, 0, 0, 0, 0, time.UTC)
	expectedExpT = time.Date(2023, time.Month(7), 12, 0, 0, 0, 0, time.UTC)

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
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgdfgdfxdfg ,sdfzsdfzsdf,D\nDAQ\nFFGG5566\nDBA\n20230712\nDBB\n19690305\nasdsda dsadASD F\nasddsaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "20230712",
			expectedExpT:   &expectedExpT,
		},
		//  "DAQ", "DBB", and "DBA" exists in data
		{
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgd DAQfgdfxdfg ,sdfzsdfzsdf,D\nDAQ\nFFGG5566\nDBA\n20230712\nDBB\n19690305\nasds DBBda dsadASD F\nasdds DBAaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "20230712",
			expectedExpT:   &expectedExpT,
		},
		// invalid expiry date
		{
			str:            "@\n\u001C\nANSI fgdgdfgxdfggfddfgdfgdfxdfg ,sdDAQfzsdfzsdf,D\nDAQ\nFFGG5566\nDBA\nabc\nDBB\n19690305\nasdsda dsadASD F\nasddsaasda\nDAJIL\nDAK600160000  \nDARD\nDASB         \nDAT*****\nDBD12345678\nDBCM\nDAU507\nDAW150\nDAYGRN\n\n",
			expectedSerial: "FFGG5566",
			expectedDOB:    "19690305",
			expectedDOBT:   &expectedDOBT,
			expectedExp:    "",
			expectedExpT:   nil,
		},
	}
)

func TestDocumentSerialFromBarcodeData(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := drivers_license.NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		require.Equal(t, tt.expectedSerial, bc.DocumentSerial)
	}
}

func TestDOBFromBarcodeData(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := drivers_license.NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		require.Equal(t, tt.expectedDOB, bc.Dob)
		require.Equal(t, tt.expectedDOBT, bc.DobT)
	}
}
func TestExpiryFromBarcodeData(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := drivers_license.NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		require.Equal(t, tt.expectedExp, bc.Expiry)
		require.Equal(t, tt.expectedExpT, bc.ExpiryT)
	}
}

func TestCompareDate(t *testing.T) {
	for _, tt := range barcodeTests {
		bc, err := drivers_license.NewBarcode(tt.str)
		require.Nil(t, err)
		require.NotNil(t, bc)
		dob, err := bc.SelectDate(drivers_license.BarcodeDataTypeDOB, tt.expectedDOBT)
		require.Nil(t, err)
		require.Equal(t, tt.expectedDOBT, dob)
		exp, err := bc.SelectDate(drivers_license.BarcodeDataTypeExpiry, tt.expectedExpT)
		require.Nil(t, err)
		require.Equal(t, tt.expectedExpT, exp)
	}

	// should return an error
	bcdata := "invalid barcode data"
	_, err := drivers_license.NewBarcode(bcdata)
	require.NotNil(t, err)
	require.IsType(t, drivers_license.ErrInvalidData, err)

	// should return an error and the original DOB
	bcdata = "invalid\nbarcode\ndata\nthat passes inspection"
	bc, err := drivers_license.NewBarcode(bcdata)
	require.NotNil(t, err)
	require.NotNil(t, bc)
	dob, err := bc.SelectDate(drivers_license.BarcodeDataTypeDOB, barcodeTests[0].expectedDOBT)
	require.Nil(t, err)
	require.Equal(t, barcodeTests[0].expectedDOBT, dob)
	exp, err := bc.SelectDate(drivers_license.BarcodeDataTypeExpiry, barcodeTests[0].expectedExpT)
	require.Nil(t, err)
	require.Equal(t, barcodeTests[0].expectedExpT, exp)

	bc, err = drivers_license.NewBarcode(barcodeTests[0].str)
	require.Nil(t, err)
	require.NotNil(t, bc)

	// should return a ErrBarcodeDateMismatch error and the DOB in the barcode
	passedDOBT := time.Date(2011, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	dob, err = bc.SelectDate(drivers_license.BarcodeDataTypeDOB, &passedDOBT)
	require.NotNil(t, err)
	require.IsType(t, drivers_license.ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedDOBT, dob)

	var errDateMismatch drivers_license.ErrBarcodeDateMismatch
	if !errors.As(err, &errDateMismatch) {
		t.Fatal("invalid error type using errors.As")
	}

	// should return a ErrBarcodeDateMismatch error and the expiration date in the barcode
	passedExpT := time.Date(2011, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	dob, err = bc.SelectDate(drivers_license.BarcodeDataTypeExpiry, &passedExpT)
	require.NotNil(t, err)
	require.IsType(t, drivers_license.ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedExpT, dob)

	if !errors.As(err, &errDateMismatch) {
		t.Fatal("invalid error type using errors.As")
	}

	// should return a ErrBarcodeDateMismatch error and the expiration date in the barcode
	passedExpT = time.Time{}
	dob, err = bc.SelectDate(drivers_license.BarcodeDataTypeExpiry, &passedExpT)
	require.NotNil(t, err)
	require.IsType(t, drivers_license.ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedExpT, dob)

	if !errors.As(err, &errDateMismatch) {
		t.Fatal("invalid error type using errors.As")
	}

	// should return a ErrBarcodeDateMismatch error and the expiration date in the barcode
	dob, err = bc.SelectDate(drivers_license.BarcodeDataTypeExpiry, nil)
	require.NotNil(t, err)
	require.IsType(t, drivers_license.ErrBarcodeDateMismatch{}, err)
	require.Equal(t, barcodeTests[0].expectedExpT, dob)

	if !errors.As(err, &errDateMismatch) {
		t.Fatal("invalid error type using errors.As")
	}
}
