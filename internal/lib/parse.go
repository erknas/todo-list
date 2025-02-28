package lib

import (
	"fmt"
	"strconv"
)

func ParseID(id string) (int, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return -1, fmt.Errorf("invalid ID")
	}

	if parsedID <= 0 {
		return -1, fmt.Errorf("invalid ID")
	}

	return parsedID, nil
}
