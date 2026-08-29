package service

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

var derivedNameTokenPattern = regexp.MustCompile(`\{([^{}]+)\}`)

func validateDerivedNameFormat(
	ctx context.Context,
	q repository.Querier,
	format string,
	properties []dto.ItemTypeProperty,
) error {
	tokens, valid := derivedNameTokens(format)
	if !valid || len(tokens) == 0 {
		return ErrInvalidDerivedNameFormat
	}

	propertyNames := make(map[string]struct{}, len(properties))
	seenIDs := make(map[int64]struct{}, len(properties))
	for _, property := range properties {
		if _, seen := seenIDs[property.ID]; seen {
			continue
		}
		seenIDs[property.ID] = struct{}{}

		itemProperty, err := q.GetPropertyByID(ctx, property.ID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return ErrInvalidReference
			}
			return err
		}
		propertyNames[itemProperty.Name] = struct{}{}
	}

	for _, token := range tokens {
		if _, exists := propertyNames[token]; !exists {
			return ErrInvalidDerivedNameFormat
		}
	}

	return nil
}

func derivedNameTokens(format string) ([]string, bool) {
	if strings.TrimSpace(format) == "" {
		return nil, false
	}

	matches := derivedNameTokenPattern.FindAllStringSubmatchIndex(format, -1)
	if len(matches) == 0 {
		return nil, false
	}

	tokens := make([]string, 0, len(matches))
	lastEnd := 0
	for _, match := range matches {
		if strings.ContainsAny(format[lastEnd:match[0]], "{}") {
			return nil, false
		}

		token := strings.TrimSpace(format[match[2]:match[3]])
		if token == "" {
			return nil, false
		}
		tokens = append(tokens, token)
		lastEnd = match[1]
	}

	if strings.ContainsAny(format[lastEnd:], "{}") {
		return nil, false
	}

	return tokens, true
}

func derivedNameUsesProperty(format, propertyName string) bool {
	for _, token := range derivedNameTokenPattern.FindAllStringSubmatch(format, -1) {
		if strings.TrimSpace(token[1]) == propertyName {
			return true
		}
	}
	return false
}

func populateItemDetails(
	items []dto.Item,
	propertyRows []repository.GetItemPropertiesRow,
	derivedNameRows []repository.GetItemsDerivedNamesRow,
) {
	itemIndexes := make(map[int64]int, len(items))
	for index := range items {
		itemIndexes[items[index].ID] = index
	}

	for _, row := range propertyRows {
		itemIndex, exists := itemIndexes[row.ItemProperty.ItemID]
		if !exists {
			continue
		}

		property := dto.ToItemPropertyDTO(row.ItemProperty)
		property.Visibility = string(row.Visibility)
		property.ValueType = row.PropertyType
		switch property.ValueType {
		case "expiry":
			var propertyValue string
			if err := json.Unmarshal(property.Value, &propertyValue); err != nil {
				log.Printf("Couldn't unmarshal expiry date value %v. Error: %v", string(property.Value), err)
				continue
			}
			expiryTime, err := time.Parse(time.DateOnly, propertyValue)
			if err != nil {
				log.Printf("Couldn't parse expiry date %v. Error: %v", propertyValue, err)
				continue
			}
			currentTime := time.Now().UTC()
			if currentTime.After(expiryTime) {
				property.SmartData = "expired"
			} else {
				difference := expiryTime.Sub(currentTime)
				days := difference.Hours() / 24
				if days < 7 {
					property.SmartData = "expiring_soon"
				}
			}
		}
		items[itemIndex].Properties = append(items[itemIndex].Properties, property)
	}

	for _, row := range derivedNameRows {
		itemIndex, exists := itemIndexes[row.ItemID]
		if !exists {
			continue
		}

		items[itemIndex].DerivedName = row.DerivedName
	}
}
