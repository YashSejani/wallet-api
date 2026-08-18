package api

import (
	"strings"
)

const (
	USD = "USD"
	EUR = "EUR"
	INR = "INR"
	GBP = "GBP"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch strings.ToUpper(currency) {
	case USD, EUR, INR, GBP, CAD:
		return true
	}
	return false
}

func errorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
