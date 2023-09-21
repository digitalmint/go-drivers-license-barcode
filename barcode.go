package drivers_license

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
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

type DateField struct {
	String string
	DateT  *time.Time
	Err    error
}
type StringField struct {
	String string
	Err    error
}
type Barcode struct {
	Raw            string
	Dob            DateField
	Expiry         DateField
	DocumentSerial StringField
}

// NewBarcode instantiates and returns a Barcode object.
// If the barcode data is invalid and cannot be parsed, it will return an empty Barcode and the error.
// Otherwise, it attempts to parse the serial number, dob and expiration date
// Any errors parsing these are stored in the respective Barcode field's Err value rather than failing.
func NewBarcode(data string) (Barcode, error) {
	if !strings.Contains(data, "\n") {
		return Barcode{}, ErrInvalidData{}
	}
	data = strings.TrimSpace(data)
	bc := Barcode{Raw: data}

	bc.DocumentSerial.String, bc.DocumentSerial.Err = extractData(data, BarcodeDataPrefixSerial)

	bc.Dob = processDate(data, BarcodeDataPrefixDOB)

	bc.Expiry = processDate(data, BarcodeDataPrefixExpiry)

	return bc, nil
}

// SelectDate compares the date of the type BarcodeDataType found in the barcode data with a date which is passed in time.Time format.
// If dates do not match, it returns the barcode date and a ErrBarcodeDateMismatch error, otherwise it returns the original date passed in, and any error.
// If no barcode date was found, it returns the original date passed in, and nil error.
// Dates are returned in time.Time format
func (bc Barcode) SelectDate(dateType BarcodeDataType, date *time.Time) (*time.Time, error) {
	var bcDate string

	switch dateType {
	case BarcodeDataTypeDOB:
		bcDate = bc.Dob.String

	case BarcodeDataTypeExpiry:
		bcDate = bc.Expiry.String

	default:
		zap.S().Panicf("invalid dateType: %s", dateType)
	}

	fieldName := string(dateType)

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
			return date, ErrParseDate{
				Date:      bcDate,
				FieldName: fieldName,
				Err:       err,
			}
		}
		date = &tmp
		return date, ErrBarcodeDateMismatch{
			SentDate:    sentDateStr,
			BarcodeDate: bcDate,
			FieldName:   fieldName,
		}
	}
	return date, nil
}

// processDate extracts the date for the prefix, and return it in *time.Time and as a YYYYMMDD formatted string, and any errors
func processDate(data string, prefix BarcodeDataPrefix) DateField {
	date, err := extractData(data, prefix)
	if err != nil {
		return DateField{
			Err: err,
		}
	}
	if date == "" {
		return DateField{}
	}

	fieldName := "uknown"
	switch prefix {
	case BarcodeDataPrefixDOB:
		fieldName = string(BarcodeDataTypeDOB)

	case BarcodeDataPrefixExpiry:
		fieldName = string(BarcodeDataTypeExpiry)

	}

	dateT, err := parseDate(date, fieldName)

	if err != nil {
		return DateField{
			Err: err,
		}
	}
	// convert to a YYYYMMDD standard format
	date = dateT.Format(TimeLayoutBarcodeData)

	return DateField{
		String: date,
		DateT:  dateT,
		Err:    nil,
	}
}

// parseDate takes in a date and field name strings
// It will then attempt to convert to a time.Time value and return it with any errors
func parseDate(date, fieldName string) (*time.Time, error) {
	_, err := strconv.Atoi(date)
	if err != nil {
		return nil, ErrInvalidDate{
			FieldName: fieldName,
			Value:     date,
		}
	}
	yy, _ := strconv.Atoi(date[0:2])

	var dateT time.Time
	if yy >= 19 {
		//YYYYMMDD
		dateT, err = time.Parse(TimeLayoutBarcodeData, date)
	} else {
		//MMDDYYYY
		dateT, err = time.Parse(TimeLayoutBarcodeDataUS, date)
	}
	if err != nil {
		return nil, ErrParseDate{
			Date:      date,
			FieldName: fieldName,
			Err:       err,
		}
	}
	return &dateT, nil
}

// extractData extracts the information associated with 'prefix' from 'data' and return the extracted string and any errors
func extractData(data string, prefix BarcodeDataPrefix) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`\n%s\s*(\S+)`, prefix))
	match := re.FindStringSubmatch(data)
	if len(match) > 1 {
		return strings.TrimSpace(match[1]), nil
	} else {
		isDate := prefix == BarcodeDataPrefixDOB || prefix == BarcodeDataPrefixExpiry
		return "", ErrPrefixExtraction{
			Prefix:      prefix,
			IsDateError: isDate,
		}
	}

}
