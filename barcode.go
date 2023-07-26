package drivers_license

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidData = errors.New("invalid barcode data sent")
	ErrInvalidDate = errors.New("invalid date")
)

type BarcodeDataType string
type BarcodeDataPrefix string

const (
	BarcodeDataTypeDOB    BarcodeDataType = "dob"
	BarcodeDataTypeExpiry BarcodeDataType = "exp"

	BarcodeDataPrefixSerial BarcodeDataPrefix = "DAQ"
	BarcodeDataPrefixDOB    BarcodeDataPrefix = "DBB"
	BarcodeDataPrefixExpiry BarcodeDataPrefix = "DBA"

	TimeLayoutBarcodeData   = "20060102"
	TimeLayoutBarcodeDataUS = "01022006"
)

type ErrBarcodeDateMismatch struct {
	SentDOB, BarcodeDOB string
}

func (e ErrBarcodeDateMismatch) Error() string {
	return fmt.Sprintf("barcode date (%s) does not match passed date (%s) - using barcode date", e.BarcodeDOB, e.SentDOB)
}

type Barcode struct {
	Raw            string
	Dob            string
	DobT           *time.Time
	Expiry         string
	ExpiryT        *time.Time
	DocumentSerial string
}

func NewBarcode(data string) (Barcode, error) {
	var err error
	if !strings.Contains(data, "\n") {
		return Barcode{}, ErrInvalidData
	}
	data = strings.TrimSpace(data)
	bc := Barcode{Raw: data}

	bc.DocumentSerial, err = extractData(data, BarcodeDataPrefixSerial, nil)
	if err != nil {
		return bc, err
	}

	bc.DobT, bc.Dob, err = processDate(data, BarcodeDataPrefixDOB)
	if err != nil && !errors.Is(err, ErrInvalidDate) {
		return bc, err
	}

	bc.ExpiryT, bc.Expiry, err = processDate(data, BarcodeDataPrefixExpiry)
	if err != nil && !errors.Is(err, ErrInvalidDate) {
		return bc, err
	}

	return bc, nil
}

// SelectDate compares the date of the type dateType found in the barcode data with a date which is passed in time.Time format.
// If dates do not match, it returns the barcode date and a ErrBarcodeDateMismatch error, otherwise it returns the original date passed in, and any error.
// If no barcode date was found it returns the original date passed in, and any error.
// Dates are returned in time.Time format
func (bc Barcode) SelectDate(dateType BarcodeDataType, date *time.Time) (*time.Time, error) {
	var bcDate string

	switch dateType {
	case BarcodeDataTypeDOB:
		bcDate = bc.Dob

	case BarcodeDataTypeExpiry:
		bcDate = bc.Expiry

	default:
		zap.S().Panicf("invalid dateType: %s", dateType)
	}

	// check for empty values
	if bcDate == "" {
		return date, nil
	}
	if date == nil {
		date = &time.Time{}
	}
	sentDateStr := date.Format(TimeLayoutBarcodeData)
	if sentDateStr != bcDate {
		tmp, err := time.Parse(TimeLayoutBarcodeData, bcDate)
		if err != nil {
			return date, err
		}
		date = &tmp
		return date, ErrBarcodeDateMismatch{
			SentDOB:    sentDateStr,
			BarcodeDOB: bcDate,
		}
	}
	return date, nil
}

func processDate(data string, prefix BarcodeDataPrefix) (*time.Time, string, error) {
	date, err := extractData(data, prefix, nil)
	if err != nil {
		return nil, "", err
	}
	if date == "" {
		return nil, "", nil
	}
	_, err = strconv.Atoi(date)
	if err != nil {
		return nil, "", ErrInvalidDate
	}

	yy, _ := strconv.Atoi(date[0:2])

	var dateT time.Time
	if yy >= 19 {
		//YYYYMMDD
		dateT, err = time.Parse(TimeLayoutBarcodeData, date)
	} else {
		//MMDDYYYY
		dateT, err = time.Parse(TimeLayoutBarcodeDataUS, date)
		// convert to a YYYYMMDD standard format
		date = dateT.Format(TimeLayoutBarcodeData)
	}
	if err != nil {
		return nil, "", err
	}

	return &dateT, date, nil
}

func extractData(data string, prefix BarcodeDataPrefix, retErr error) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`\n%s\s*(\S+)`, prefix))
	match := re.FindStringSubmatch(data)
	if retErr == nil {
		retErr = fmt.Errorf("barcode data using prefix: %s could not be extracted from the barcode data", prefix)
	}
	if len(match) > 1 {
		return strings.TrimSpace(match[1]), nil
	} else {
		return "", retErr
	}

}
