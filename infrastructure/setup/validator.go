package setup

import (
"github.com/go-playground/validator/v10"
"github.com/spf13/viper"
)

// NewValidator creates a new validator instance
func NewValidator(viper *viper.Viper) *validator.Validate {
	return validator.New()
}
