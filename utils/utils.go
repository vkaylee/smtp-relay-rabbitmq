package utils

import "log"

/*
CheckErr checks for error, and log fatal if it is not nil.
*/
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
