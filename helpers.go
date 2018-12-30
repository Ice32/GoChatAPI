package main

import "fmt"

func sliceContainsString(slice []string, value string) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	return false
}

func logError(error string) {
	fmt.Println(error)
}
