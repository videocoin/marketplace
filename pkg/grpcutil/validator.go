package grpcutil

import (
	"reflect"
	"strings"
	"unicode"

	enLocale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/videocoin/marketplace/api/rpc"
	validator "gopkg.in/go-playground/validator.v9"
	enTrans "gopkg.in/go-playground/validator.v9/translations/en"
)

type RequestValidator struct {
	validator  *validator.Validate
	translator *ut.Translator
}

func NewRequestValidator() (*RequestValidator, error) {
	lt := enLocale.New()
	en := &lt

	uniTranslator := ut.New(*en, *en)
	uniEn, _ := uniTranslator.GetTranslator("en")
	translator := &uniEn

	validate := validator.New()
	err := enTrans.RegisterDefaultTranslations(validate, *translator)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterTranslation(
		"email",
		*translator,
		RegisterEmailTranslation,
		EmailTranslation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterTranslation(
		"secure-password",
		*translator,
		RegisterSecurePasswordTranslation,
		SecurePasswordTranslation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterTranslation(
		"confirm-password",
		*translator,
		RegisterConfirmPasswordTranslation,
		ConfirmPasswordTranslation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("confirm-password", ValidateConfirmPassword)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("secure-password", ValidateSecurePassword)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterTranslation(
		"address",
		*translator,
		RegisterAddressTranslation,
		AddressTranslation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterTranslation(
		"pin",
		*translator,
		RegisterPinTranslation,
		PinTranslation)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterTranslation(
		"transfer_id",
		*translator,
		RegisterTransferIDTranslation,
		TransferIDTranslation)
	if err != nil {
		return nil, err
	}
	return &RequestValidator{
		validator:  validate,
		translator: translator,
	}, nil

}

func (rv *RequestValidator) Validate(r interface{}) *rpc.MultiValidationError {
	trans := *rv.translator
	verrs := &rpc.MultiValidationError{}

	serr := rv.validator.Struct(r)
	if serr != nil {
		verrs.Errors = []*rpc.ValidationError{}

		for _, err := range serr.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(r).Elem().FieldByName(err.Field())
			jsonField := extractValueFromTag(field.Tag.Get("json"))
			verr := &rpc.ValidationError{
				Field:   jsonField,
				Message: err.Translate(trans),
			}
			verrs.Errors = append(verrs.Errors, verr)
		}

		return verrs
	}

	return nil
}

func RegisterEmailTranslation(ut ut.Translator) error {
	return ut.Add("email", "Enter a valid email address", true)
}

func EmailTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("email", fe.Field())
	return t
}

func RegisterSecurePasswordTranslation(ut ut.Translator) error {
	return ut.Add("secure-password", "Password must be more than 8 characters and contain both numbers and letters", true)
}

func SecurePasswordTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("secure-password", fe.Field())
	return t
}

func RegisterConfirmPasswordTranslation(ut ut.Translator) error {
	return ut.Add("confirm-password", "Passwords does not match", true)
}

func ConfirmPasswordTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("confirm-password", fe.Field())
	return t
}

func ValidateConfirmPassword(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	return field.String() == currentField.String()
}

func ValidateSecurePassword(fl validator.FieldLevel) bool {
	field := fl.Field()
	password := field.String()

	if password == "" {
		return false
	}

	var (
		hasMinLen = false
		hasNumber = false
		hasLetter = false
		// hasUpper   = false
		// hasSpecial = false
		// hasLower   = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char):
			hasNumber = true
			//     hasUpper = true
			// case unicode.IsLower(char):
			//     hasLower = true
			// case unicode.IsPunct(char) || unicode.IsSymbol(char):
			//     hasSpecial = true
		}
	}

	return hasMinLen && hasNumber && hasLetter
}

func RegisterAddressTranslation(ut ut.Translator) error {
	return ut.Add("address", "Enter a valid ethereum address", true)
}

func AddressTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("address", fe.Field())
	return t
}

func RegisterPinTranslation(ut ut.Translator) error {
	return ut.Add("pin", "Enter a valid pin", true)
}

func PinTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("pin", fe.Field())
	return t
}

func RegisterTransferIDTranslation(ut ut.Translator) error {
	return ut.Add("transfer_id", "Enter a valid transfer id", true)
}

func TransferIDTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("transfer_id", fe.Field())
	return t
}

func extractValueFromTag(tag string) string {
	values := strings.Split(tag, ",")
	return values[0]
}
