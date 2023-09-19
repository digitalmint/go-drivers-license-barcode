package drivers_license

import (
	"errors"
	"fmt"
)

type ErrInvalidData struct{}

func (e ErrInvalidData) Error() string {
	return "invalid barcode data"
}

type ErrInvalidDate struct{}

func (e ErrInvalidDate) Error() string {
	return "invalid date"
}

type ErrBarcodeDateMismatch struct {
	SentDOB, BarcodeDOB string
}

func (e ErrBarcodeDateMismatch) Error() string {
	return fmt.Sprintf("barcode date %q does not match passed date %q - using barcode date", e.BarcodeDOB, e.SentDOB)
}

type ErrPrefixExtraction struct {
	Prefix BarcodeDataPrefix
}

func (e ErrPrefixExtraction) Error() string {
	return fmt.Sprintf("prefix: %q could not be extracted from the barcode data", e.Prefix)
}

type ErrParseDate struct {
	Date string
	Err  error
}

func (e ErrParseDate) Error() string {
	return fmt.Sprintf("could not parse date %q", e.Date)
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
