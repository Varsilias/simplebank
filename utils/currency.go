package utils

// Constants for all supported currencies
const (
	NGN = "NGN"
	EUR = "EUR"
	USD = "USD"
	GBP = "GBP"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case NGN, EUR, USD, GBP:
		return true
	}
	return false
}
