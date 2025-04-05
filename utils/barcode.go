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

	sumOdd := 0  // positions 1,3,5,7,9,11 (0-based index)
	sumEven := 0 // positions 2,4,6,8,10,12

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
