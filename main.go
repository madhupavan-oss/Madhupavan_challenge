package main

import (
	"fmt"
	"regexp"
)

func isValidCreditCard(cardNumber string) bool {

	pattern := `^[456]\d{3}(-?\d{4}){3}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(cardNumber) {
		return false
	}
	for i := 0; i < len(cardNumber)-3; i++ {
		if cardNumber[i] == cardNumber[i+1] && cardNumber[i+1] == cardNumber[i+2] && cardNumber[i+2] == cardNumber[i+3] {
			return false
		}
	}

	return true
}

func main() {
	// Test cases
	cardNumbers := []string{
		"4253625879615786",
		"4424424424442444",
		"5122-2368-7954-3214",

		"42536258796157867",
		"4424444424442444",
		"5122-2368-7954 - 3214",
		"44244x4424442444",
		"0525362587961578 ",
	}

	// Check validity for each card number
	for _, cardNumber := range cardNumbers {
		if isValidCreditCard(cardNumber) {
			fmt.Println(cardNumber, "is valid")
		} else {
			fmt.Println(cardNumber, "is not valid")
		}
	}
}
