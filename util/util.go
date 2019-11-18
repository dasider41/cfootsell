package util

import (
	"log"
	"regexp"
	"strconv"
)

// ErrCheck :
func ErrCheck(err error) {
	if err != nil {
		log.Fatal(err)
		// TODO :: Report an error by eamil
	}
}

// NumberOnly :
func NumberOnly(text string) (int, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(reg.ReplaceAllString(text, ""))
	if err != nil {
		return 0, err
	}
	return val, nil
}
