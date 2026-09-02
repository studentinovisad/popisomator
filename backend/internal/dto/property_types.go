package dto

import (
	"cmp"
	"maps"
	"slices"
)

const PriceMultiplier = 10000

// Property Type - Price
type PTPrice struct {
	Amount   int64  `json:"amount" validate:"gte=0"`
	Currency string `json:"currency" validate:"iso4217"`
}

// MeasureMultiplier scales the amount of a measured property type (mass, volume) the same way
// PriceMultiplier scales a price: the amount is stored as an integer in the unit the user picked,
// multiplied by this, so four decimal places survive without floating point drift.
const MeasureMultiplier = 10000

// Property Type - Mass
type PTMass struct {
	Amount int64  `json:"amount" validate:"gte=0"`
	Unit   string `json:"unit" validate:"oneof=mg g kg"`
}

// Property Type - Volume
type PTVolume struct {
	Amount int64  `json:"amount" validate:"gte=0"`
	Unit   string `json:"unit" validate:"oneof=mL L"`
}

// MassUnitFactors and VolumeUnitFactors convert an amount into the canonical base unit for its
// dimension - milligram for mass, millilitre for volume - so amounts recorded in different units
// can be summed: baseAmount = amount * factor, still scaled by MeasureMultiplier. Keep the keys in
// sync with the oneof tags on PTMass.Unit and PTVolume.Unit.
var (
	MassUnitFactors   = map[string]int64{"mg": 1, "g": 1_000, "kg": 1_000_000}
	VolumeUnitFactors = map[string]int64{"mL": 1, "L": 1_000}
)

// MeasureUnitFactorRows flattens MassUnitFactors and VolumeUnitFactors into three parallel arrays,
// one row per (value type, unit, factor), for SumItemProperties to join against. Passing the table
// in as query parameters keeps it defined only here instead of also being spelled out in SQL. The
// rows are sorted so the parameters a query gets are stable.
func MeasureUnitFactorRows() (valueTypes []string, units []string, factors []int64) {
	factorsByValueType := map[string]map[string]int64{
		"mass":   MassUnitFactors,
		"volume": VolumeUnitFactors,
	}

	for _, valueType := range slices.Sorted(maps.Keys(factorsByValueType)) {
		unitFactors := factorsByValueType[valueType]
		for _, unit := range slices.Sorted(maps.Keys(unitFactors)) {
			valueTypes = append(valueTypes, valueType)
			units = append(units, unit)
			factors = append(factors, unitFactors[unit])
		}
	}

	return valueTypes, units, factors
}

// ScaleMeasureToLargestUnit re-expresses an amount given in a dimension's base unit as the largest
// unit that still leaves at least one of it, so a summed mass reads 3.456 kg rather than 3456000 mg.
// It only steps up while the conversion divides exactly, so a total is never rounded away: a mass
// carrying detail finer than 0.1 g reports in grams instead of kilograms. The base unit - the one
// with factor 1 - always divides exactly, so there is always an answer.
func ScaleMeasureToLargestUnit(baseAmount int64, unitFactors map[string]int64) (int64, string) {
	largestFirst := slices.SortedFunc(maps.Keys(unitFactors), func(left, right string) int {
		return cmp.Compare(unitFactors[right], unitFactors[left])
	})

	for _, unit := range largestFirst {
		factor := unitFactors[unit]
		if baseAmount >= MeasureMultiplier*factor && baseAmount%factor == 0 {
			return baseAmount / factor, unit
		}
	}

	baseUnit := largestFirst[len(largestFirst)-1]
	return baseAmount / unitFactors[baseUnit], baseUnit
}
