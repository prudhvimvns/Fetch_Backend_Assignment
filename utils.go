package main

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func CalculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: Alphanumeric characters in retailer name
	for _, c := range receipt.Retailer {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			points++
		}
	}

	// Convert total to float64 once
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return points // fallback to whatever points accumulated
	}

	// Rule 2: Round dollar (no cents)
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule 3: Total is a multiple of 0.25
	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Trimmed description length % 3
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 6: Large Language Model bonus if total > 10.00
	if totalFloat > 10.00 {
		points += 5
	}

	// Rule 7: Odd day of purchase
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if date.Day()%2 == 1 {
			points += 6
		}
	}

	// Rule 8: Purchase time between 2:00PM and 4:00PM
	if t, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if t.Hour() == 14 {
			points += 10
		}
	}

	return points
}
