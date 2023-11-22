package api

import (
	"github.com/go-playground/validator/v10"
	"simplebank/utils"
)

var validCureency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// eheck currency is supported
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
