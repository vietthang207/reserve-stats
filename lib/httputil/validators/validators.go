package validators

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gopkg.in/go-playground/validator.v8"

	"reflect"
)

// isEthereumAddress is a validator.Func function that returns true if given field
// is a valid Ethereum address.
func isEthereumAddress(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	address := field.String()
	if err := validation.Validate(address, validation.Required); err != nil {
		return false
	}
	if !common.IsHexAddress(address) {
		return false
	}
	return true
}

// isEmail is a validator.Func function that returns true if given field
// is a valid email address.
func isEmail(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	if err := validation.Validate(field.String(), is.Email); err != nil {
		return false
	}
	return true
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isAddress", isEthereumAddress)
		v.RegisterValidation("isEmail", isEmail)
	}
}
