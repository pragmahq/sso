/*
			    Date - Some library code to handle just "Dates" because I don't want to complicate using go timestamps.

	    Sample usage:
		            date := internal.NewDate(6, 12, 1990)
		            fmt.Printf("Created Date: %v\n", date)

		            dateString := "15/03/1985"
		            if err := internal.ParseDateFromString(dateString, &date); err != nil {
		                fmt.Printf("Error parsing Date: %v\n", err)
		            } else {
		                fmt.Printf("Parsed Date: %v\n", date)
		            }

		            fmt.Printf("Date as string: %s\n", internal.DateToString(date))
*/
package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Date struct {
	Day   int
	Month int
	Year  int
}

// Return an instance of Date from DDMMYYYY integers
func NewDate(day, month, year int) Date {
	return Date{
		Day:   day,
		Month: month,
		Year:  year,
	}
}

// ParseDateFromString parses a string in the format "DD/MM/YYYY" and sets the values in the Date instance.
func ParseDateFromString(dateString string, d *Date) error {
	parts := strings.Split(dateString, "/")

	if len(parts) != 3 {
		return fmt.Errorf("invalid date format: %s", dateString)
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("failed to parse day: %v", err)
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse month: %v", err)
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse year: %v", err)
	}

	d.Day = day
	d.Month = month
	d.Year = year

	return nil
}

// DateToString returns a string representation of the Date in the format "DD/MM/YYYY".
func DateToString(d Date) string {
	return fmt.Sprintf("%02d/%02d/%04d", d.Day, d.Month, d.Year)
}
