package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vldtr "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidatorImpl struct {
	validator *vldtr.Validate
	trans     ut.Translator
}

func NewValidator() *ValidatorImpl {
	validator := vldtr.New()
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	// Register default english translations
	en_translations.RegisterDefaultTranslations(validator, trans)

	return &ValidatorImpl{
		validator: validator,
		trans:     trans,
	}
}

// Validate function validates the input data
func (cv *ValidatorImpl) Validate(data interface{}) error {

	var msgs string
	err := cv.validator.Struct(data)
	if err != nil {
		for idx, e := range err.(vldtr.ValidationErrors) {
			if idx > 0 && idx < len(err.(vldtr.ValidationErrors)) {
				msgs = msgs + ", "
			}

			msgs = msgs + strings.ToLower(e.Translate(cv.trans))
		}
		return errors.New(msgs)
	}
	return nil
}
