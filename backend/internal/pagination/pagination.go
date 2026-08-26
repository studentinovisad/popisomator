package pagination

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultPageSize     int32 = 20
	MinimumPageSize     int32 = 1
	MaximumPageSize     int32 = 50
	MaximumOffset       int32 = 1_000_000
	MaximumSearchLength       = 100
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
	if parsedValue < minimum || parsedValue > maximum {
		return 0, strconv.ErrSyntax
	}

	return parsedValue, nil
}

func GetLimitOffset(request *http.Request) (int32, int32, error) {
	limit, err := GetLimit(request)
	if err != nil {
		return 0, 0, err
	}

	offset, err := queryValue(request, "offset", 0, 0, MaximumOffset)
	if err != nil {
		return 0, 0, err
	}

	return limit, offset, nil
}

func GetLimit(request *http.Request) (int32, error) {
	return queryValue(request, "limit", DefaultPageSize, MinimumPageSize, MaximumPageSize)
}

func GetSearch(request *http.Request) (string, error) {
	search := strings.TrimSpace(request.URL.Query().Get("search"))
	if len(search) > MaximumSearchLength {
		return "", errors.New("search is too long")
	}

	return search, nil
}
