package handler

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/AthanatiusC/SawitPro/generated"
)

/*
- Validate interface user scheme using reflect to get value and types
- Reflect instead of custom validator because this offers more flexibility
- Iterate and validate each field by tag, complexity O(n)
*/
func (s *Server) ValidateUser(request interface{}) (errors generated.ErrorValidationResponse) {
	v := reflect.ValueOf(request)
	validations := make(map[string][]string)

	for i := 0; i < v.NumField(); i++ {
		var lowerLimit, upperLimit int
		tag := v.Type().Field(i).Tag.Get("json") // Use json tag only because register use request body
		val := v.Field(i).Interface()

		const (
			FullNameTag    = "full_name"
			PasswordTag    = "password"
			PhoneNumberTag = "phone_number"
		)

		switch tag {
		case FullNameTag:
			lowerLimit, upperLimit = 3, 60
		case PasswordTag:
			upperLimit, lowerLimit = 64, 6
			password := v.Field(i).String()
			// Golang regexp does not support lookaround, spearate each regex instead of one combination
			specialRegexp := regexp.MustCompile(`[!@#$%^&*()_+\[\]{};':"\|,.<>?]`)
			capitialRegexp := regexp.MustCompile(`[A-Z]`)
			numericalRegexp := regexp.MustCompile(`[0-9]`)
			if !specialRegexp.MatchString(password) {
				validations[tag] = append(validations[tag], "at least 1 special character")
			}
			if !capitialRegexp.MatchString(password) {
				validations[tag] = append(validations[tag], "at least 1 capital character")
			}
			if !numericalRegexp.MatchString(password) {
				validations[tag] = append(validations[tag], "at least 1 number")
			}
		case PhoneNumberTag:
			/*
				Flow:
				1. Validate phone number string to given regex pattern
				2. Remove non-alphanumeric characters for utf-8 to ensure only numeric is phone number
				3. Validate if phone number length is more than character limit
			*/
			lowerLimit, upperLimit = 0, 13
			pattern, _ := regexp.Compile(`^(^\+62|62)(\d{3,4}-?){2}\d{3,4}$`) // Ignore error because value is hardcoded
			phoneNumber := pattern.FindAllString(v.Field(i).String(), -1)
			if len(phoneNumber) == 0 {
				validations[tag] = append(validations[tag], "must be indonesian(+62) format ")
			}
			validPhoneNumber := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(v.Field(i).String(), "")
			if len(validPhoneNumber) < lowerLimit || len(validPhoneNumber) > upperLimit {
				validations[tag] = append(validations[tag], fmt.Sprintf("must be more than %d and less than %d characters long", lowerLimit, upperLimit))
			}
		}

		// Default validation for every field
		if val == "" {
			validations[tag] = append([]string{}, fmt.Sprintf("%s is required", tag))
		} else if (upperLimit != 0 && lowerLimit != 0) && v.Field(i).Len() < lowerLimit || v.Field(i).Len() > upperLimit && tag != PhoneNumberTag { // Validate length except phone number
			validations[tag] = append(validations[tag], fmt.Sprintf("must be more than %d and less than %d characters long", lowerLimit, upperLimit))
		}
	}

	// If error exist, construct, format and return error message
	if len(validations) != 0 {
		for field, validation := range validations {
			errors.Messages = append(errors.Messages, fmt.Sprintf("%s : %s", field, strings.Join(validation, ", ")))
		}
		return errors
	}
	return
}

func CleanPhoneNumber(phoneNumber string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(phoneNumber, "")
}
