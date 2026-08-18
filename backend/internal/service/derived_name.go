package service

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

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

func renderDerivedName(format string, values map[string]string) string {
	if _, valid := derivedNameTokens(format); !valid {
		return ""
	}

	return strings.TrimSpace(derivedNameTokenPattern.ReplaceAllStringFunc(format, func(token string) string {
		name := strings.TrimSpace(token[1 : len(token)-1])
		return values[name]
	}))
}

func derivedNameUsesProperty(format, propertyName string) bool {
	for _, token := range derivedNameTokenPattern.FindAllStringSubmatch(format, -1) {
		if strings.TrimSpace(token[1]) == propertyName {
			return true
		}
	}
	return false
}

func displayDerivedNameValue(rawValue string) string {
	var value any
	if err := json.Unmarshal([]byte(rawValue), &value); err != nil {
		return rawValue
	}

	switch typedValue := value.(type) {
	case string:
		return typedValue
	case nil:
		return ""
	default:
		encoded, err := json.Marshal(typedValue)
		if err != nil {
			return rawValue
		}
		return string(encoded)
	}
}

func populateItemDetails(items []dto.Item, rows []repository.GetItemPropertiesRow) {
	itemIndexes := make(map[int64]int, len(items))
	propertyValues := make(map[int64]map[string]string, len(items))
	formats := make(map[int64]string, len(items))
	for index := range items {
		itemIndexes[items[index].ID] = index
		propertyValues[items[index].ID] = make(map[string]string)
	}

	for _, row := range rows {
		itemIndex, exists := itemIndexes[row.ItemProperty.ItemID]
		if !exists {
			continue
		}

		property := dto.ToItemPropertyDTO(row.ItemProperty)
		property.Visibility = string(row.Visibility)
		items[itemIndex].Properties = append(items[itemIndex].Properties, property)
		propertyValues[row.ItemProperty.ItemID][row.PropertyName] = displayDerivedNameValue(
			row.ItemProperty.PropertyValue,
		)

		if row.DerivedNameFormat.Valid {
			formats[row.ItemProperty.ItemID] = row.DerivedNameFormat.String
		}
	}

	for index := range items {
		items[index].DerivedName = renderDerivedName(
			formats[items[index].ID],
			propertyValues[items[index].ID],
		)
	}
}
