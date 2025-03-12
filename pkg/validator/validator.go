package validator

import (
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	vldtr "github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

// Validator is a custom validator for the application
type Validator struct {
	validator *vldtr.Validate
	Trans     ut.Translator
}

// Validate function validates the input data
func (cv *Validator) Validate(i interface{}) error {
	idLocale := id.New()
	uni := ut.New(idLocale, idLocale)
	cv.Trans, _ = uni.GetTranslator("id")

	// Register the validator
	id_translations.RegisterDefaultTranslations(cv.validator, cv.Trans)

	return cv.validator.Struct(i)
}
