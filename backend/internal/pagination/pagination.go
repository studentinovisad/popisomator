package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPageSize int32 = 20
	MinimumPageSize int32 = 1
	MaximumPageSize int32 = 50
)

func queryValue(request *http.Request, key string, fallback, minimum, maximum int32) (int32, error) {
	value := request.URL.Query().Get(key)
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

func GetLimitOffset(request *http.Request) (int32, int32, error) {
	limit, err := queryValue(request, "limit", DefaultPageSize, MinimumPageSize, MaximumPageSize)
	if err != nil {
		return 0, 0, err
	}

	offset, err := queryValue(request, "offset", 0, 0, 0)
	if err != nil {
		return 0, 0, err
	}

	return limit, offset, nil
}
