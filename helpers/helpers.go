package helpers

import "fmt"

func SliceContainsString(slice []string, value string) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	return false
}

func LogError(error error) {
	fmt.Println(error.Error())
}
