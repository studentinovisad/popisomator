package controller

import (
	"net/http"
	"strconv"
)

const (
	defaultPageSize int32 = 20
	minimumPageSize int32 = 5
	maximumPageSize int32 = 50
)

func paginationValue(r *http.Request, key string, fallback, minimum, maximum int32) (int32, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}

	parsedValue := int32(parsed)
	if parsedValue < minimum || (maximum > 0 && parsedValue > maximum) {
		return 0, strconv.ErrSyntax
	}

	return parsedValue, nil
}
