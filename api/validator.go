package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/varsilias/simplebank/utils"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// this means that the conversion was successful, therefore we can go ahead and check of currency is valid
		return utils.IsSupportedCurrency(currency)
	}

	return false
}
