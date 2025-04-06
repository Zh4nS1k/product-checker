package utils

import "strconv"

func IsBarcodeValid(code string) bool {
	if len(code) != 13 {
		return false
	}
	digits := make([]int, 13)
	for i := 0; i < 13; i++ {
		d, err := strconv.Atoi(string(code[i]))
		if err != nil {
			return false
		}
		digits[i] = d
	}
	sumOdd, sumEven := 0, 0
	for i := 0; i < 12; i++ {
		if i%2 == 0 {
			sumOdd += digits[i]
		} else {
			sumEven += digits[i]
		}
	}
	total := sumOdd + sumEven*3
	checksum := (10 - (total % 10)) % 10
	return checksum == digits[12]
}

func GetCountryFromBarcode(code string) string {
	prefix := code[0:3]
	switch prefix {
	case "690", "691", "692", "693", "694", "695", "696", "697", "698", "699":
		return "Made in China"
	case "500", "509":
		return "Made in UK"
	case "890":
		return "Made in India"
	case "000", "001", "002", "003", "004", "005", "006", "007", "008", "009":
		return "Made in USA"
	default:
		return "Unknown origin"
	}
}
