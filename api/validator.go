package api

import (
	"strings"
)

// Supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	INR = "INR"
	GBP = "GBP"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch strings.ToUpper(currency) {
	case USD, EUR, INR, GBP, CAD:
		return true
	}
	return false
}

// errorResponse formats error messages into a standardized JSON response map
func errorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
