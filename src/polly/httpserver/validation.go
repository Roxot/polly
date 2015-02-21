package httpserver

import "unicode"

func isValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) != 10 {
		return false
	}

	for index, value := range phoneNumber {
		if index == 0 {
			if value != '0' {
				return false
			}
		} else if index == 1 {
			if value != '6' {
				return false
			}
		} else if !unicode.IsNumber(value) {
			return false
		}
	}

	return true
}
