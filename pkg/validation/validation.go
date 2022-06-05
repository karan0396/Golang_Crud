package validation

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

func Validate(p interface{})error{
	validate:=validator.New()
	
	//bounding function ValidatePassword
	//and you can join more

	validate.RegisterValidation("password",ValidatePassword)

	return validate.Struct(p)
}

//Validate password
func ValidatePassword(fl validator.FieldLevel)bool{
	var (
		upp, low, num, sym bool
		tot                uint8
	)
	
	pass:=fl.Field().String()
	
	for _, v := range pass {
		switch {
		case unicode.IsUpper(v):
			upp = true
			tot++
		case unicode.IsLower(v):
			low = true
			tot++
		case unicode.IsNumber(v):
			num = true
			tot++
		case unicode.IsPunct(v) || unicode.IsSymbol(v):
			sym = true
			tot++
		default:
			return false
		}
	}
 
	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}
 
	return true
}