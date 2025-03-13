package validator

import (
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	vldtr "github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

type ValidatorImpl struct {
	validator *vldtr.Validate
	trans     ut.Translator
}

func NewValidator() *ValidatorImpl {
	validator := vldtr.New()
	idLocale := id.New()
	uni := ut.New(idLocale, idLocale)
	trans, _ := uni.GetTranslator("id")

	// Register default indonesian translations
	id_translations.RegisterDefaultTranslations(validator, trans)

	return &ValidatorImpl{
		validator: validator,
		trans:     trans,
	}
}

// Validate function validates the input data
func (cv *ValidatorImpl) Validate(data interface{}) error {
	return cv.validator.Struct(data)
}
