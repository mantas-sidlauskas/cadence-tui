package ui

import (
	"fmt"
)

func PaddedString(value string, maxLen, pad int) string {

	valueLength := len(value)

	if maxLen-pad >= valueLength {
		return fmt.Sprintf("%*s%*s", pad, value, maxLen-valueLength-pad, " ")
	}

	newVal := value[0:maxLen-pad-3] + "..."
	return fmt.Sprintf("%*s%*s", pad, newVal, maxLen-len(newVal)-pad, " ")

}
