package validators

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

const tenMillion float64 = 10000000
const minusTenMillion float64 = -10000000

// Price is a float64 with max. two digits after comma
type Price float64

// Validate validates the price
func (m Price) Validate(formats strfmt.Registry) error {
	if err := validate.Minimum("", "body", float64(m), 0.01, false); err != nil {
		return err
	}

	if err := validate.Maximum("", "body", float64(m), tenMillion, false); err != nil {
		return err
	}

	if err := validateDigitsAfterComma("", "body", float64(m)); err != nil {
		return err
	}

	return nil
}

// LineItemPrice is a float64 with max. two digits after comma,
// can be less then zero for discounts applying
type LineItemPrice float64

// Validate validates the price
func (m LineItemPrice) Validate(formats strfmt.Registry) error {
	if m == 0 {
		return fmt.Errorf("price can not be zero")
	}

	if err := validate.Minimum("", "body", float64(m), minusTenMillion, false); err != nil {
		return err
	}

	if err := validate.Maximum("", "body", float64(m), tenMillion, false); err != nil {
		return err
	}

	if err := validateDigitsAfterComma("", "body", float64(m)); err != nil {
		return err
	}

	return nil
}

func validateDigitsAfterComma(path, in string, data float64) *errors.Validation {
	const maxDigitsAfterComma = 2

	s := strconv.FormatFloat(float64(data), 'f', -1, 64)
	if idx := strings.LastIndex(s, "."); idx >= 0 {
		if len(s)-1 > idx+maxDigitsAfterComma {
			return errors.FailedPattern(path, in, "number with max. two digits after a comma")
		}
	}

	return nil
}
