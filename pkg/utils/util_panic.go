package pkg_utils

import "fmt"

func PanicIfError(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("\n\n%s \nError: %s\n\n", message, err.Error()))
	}
}
