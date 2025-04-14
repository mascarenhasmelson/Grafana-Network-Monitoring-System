package error

import (
	"fmt"
	//  "log"
)

// IfError prints error, if it is not nil
func IfError(err error) bool {
	if err == nil {
		return false
	}

	fmt.Sprintf("%v", err)
	return true
}
