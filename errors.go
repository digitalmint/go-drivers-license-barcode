package drivers_license

import (
	"errors"
	"fmt"
)

type ErrInvalidData struct{}

func (e ErrInvalidData) Error() string {
	return "invalid barcode data"
}

type ErrInvalidDate struct {
	FieldName string
}

func (e ErrInvalidDate) Error() string {
	return fmt.Sprintf("fieldname: %q : invalid date", e.FieldName)
}

type ErrBarcodeDateMismatch struct {
	SentDate, BarcodeDate, FieldName string
}

func (e ErrBarcodeDateMismatch) Error() string {
	return fmt.Sprintf("fieldname: %q : barcode date %q does not match passed date %q - using barcode date", e.FieldName, e.BarcodeDate, e.SentDate)
}

type ErrPrefixExtraction struct {
	Prefix BarcodeDataPrefix
}

func (e ErrPrefixExtraction) Error() string {
	return fmt.Sprintf("prefix: %q could not be extracted from the barcode data", e.Prefix)
}

type ErrParseDate struct {
	Date, FieldName string
	Err             error
}

func (e ErrParseDate) Error() string {
	return fmt.Sprintf("fieldname: %q : could not parse date %q", e.FieldName, e.Date)
}
func (e ErrParseDate) Unwrap() error {
	return e.Err
}

func IsPackageError(err error) bool {
	return errors.As(err, &ErrBarcodeDateMismatch{}) ||
		errors.As(err, &ErrPrefixExtraction{}) ||
		errors.As(err, &ErrParseDate{}) ||
		errors.As(err, &ErrInvalidData{}) ||
		errors.As(err, &ErrInvalidDate{})
}
